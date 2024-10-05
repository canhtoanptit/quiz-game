package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"score-service/model"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

const (
	doneMessage = "done"

	assignorSticky     = "sticky"
	assignorRoundRobin = "roundrobin"
	assignorRange      = "range"
)

type Config struct {
	Brokers            string
	Group              string
	Topic              string
	Assignor           string
	Oldest             bool
	EnableTLS          bool
	InsecureSkipVerify bool
	ClientCertFile     string
	ClientKeyFile      string
	CACertFile         string
}

type Consumer interface {
	Consume(func(model.ConsumerMessage)) error
}

type SaramaConsumer struct {
	conf   *Config
	ready  chan bool
	client sarama.ConsumerGroup
	handle func(model.ConsumerMessage)
	ctx    context.Context
	cancel context.CancelFunc
}

func InitSaramaConsumerAndPanicIfError(conf *Config) *SaramaConsumer {
	c, err := NewSaramaConsumer(conf)
	if err != nil {
		panic(err)
	}

	return c
}

func NewSaramaConsumer(conf *Config) (*SaramaConsumer, error) {
	consumer := SaramaConsumer{
		conf:  conf,
		ready: make(chan bool),
	}

	config := getSaramaConfig(conf)

	consumer.ctx, consumer.cancel = context.WithCancel(context.Background())
	var err error
	consumer.client, err = sarama.NewConsumerGroup(strings.Split(conf.Brokers, ","), conf.Group, config)
	if err != nil {
		log.Error("[Kafka] CreatingConsumerGroupError: ", err)
		consumer.cancel()
		return nil, err
	}

	return &consumer, nil
}

func (c *SaramaConsumer) Consume(f func(model.ConsumerMessage)) error {
	log.Info("[Kafka] Starting consume message")
	c.handle = f
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.client.Consume(c.ctx, strings.Split(c.conf.Topic, ","), c); err != nil {
				log.Panic("[Kafka] ConsumeError: ", err)
			}
			if c.ctx.Err() != nil {
				log.Error("[Kafka] CtxError ", c.ctx.Err())
				return
			}
			c.ready = make(chan bool)
		}
	}()
	<-c.ready
	log.Println("[Kafka] SaramaConsumer up and running!")
	wg.Wait()
	return nil
}

func (c *SaramaConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *SaramaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *SaramaConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Infof("[Kafka] MessageClaimed, timestamp = %v, topic = %s, partition = %d, offset %d",
				message.Timestamp, message.Topic, message.Partition, message.Offset)
			session.MarkMessage(message, doneMessage)
			consumerMessage := model.ConsumerMessage{
				Key:       message.Key,
				Value:     message.Value,
				Topic:     message.Topic,
				Partition: message.Partition,
				Offset:    message.Offset,
			}
			c.handle(consumerMessage)

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func getSaramaConfig(conf *Config) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion

	defaultAssignor := []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	switch conf.Assignor {
	case assignorSticky:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case assignorRoundRobin:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case assignorRange:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		config.Consumer.Group.Rebalance.GroupStrategies = defaultAssignor
	}

	if conf.Oldest {
		log.Info("oldest")
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	if conf.EnableTLS {
		config.Net.TLS.Enable = true
		tlsConfig, err := NewTLSConfig(
			conf.InsecureSkipVerify,
			conf.ClientCertFile,
			conf.ClientKeyFile,
			conf.CACertFile,
		)
		if err != nil {
			log.Panic(context.TODO(), "load tls config error", err)
		}
		config.Net.TLS.Config = tlsConfig
	}

	return config
}

func NewTLSConfig(insecureSkipVerify bool, clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	tlsConfig.InsecureSkipVerify = insecureSkipVerify
	if !insecureSkipVerify {
		caCert, err := os.ReadFile(filepath.Clean(caCertFile))
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

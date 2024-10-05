package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"score-service/pkg/kafka"
)

type ApplicationConfig struct {
	KafkaConfig               *KafkaConfig
	UserSubmitAnswerTopic     string `envconfig:"user_submit_answer_topic" default:"user_submit_answer"`
	ScoreServiceConsumerGroup string `envconfig:"KAFKA_CONSUMER_GROUP" required:"true" default:"score-consumer-group"`
}

type KafkaConfig struct {
	Brokers            string `envconfig:"KAFKA_BROKERS" required:"true" default:"localhost:9092"`
	Assignor           string `envconfig:"KAFKA_ASSIGNOR" default:"range"`
	Oldest             bool   `envconfig:"KAFKA_OLDEST" default:"false"`
	EnableTLS          bool   `envconfig:"KAFKA_ENABLE_TLS" required:"true" default:"false"`
	InsecureSkipVerify bool   `envconfig:"KAFKA_INSECURE_SKIP_VERIFY" default:"false"`
	ClientCertFile     string `envconfig:"KAFKA_CLIENT_CERT_FILE"`
	ClientKeyFile      string `envconfig:"KAFKA_CLIENT_KEY_FILE"`
	CACertFile         string `envconfig:"KAFKA_CA_KEY"`
}

func (c *KafkaConfig) ToConfig(topic, group string) *kafka.Config {
	return &kafka.Config{
		Brokers:            c.Brokers,
		Group:              group,
		Topic:              topic,
		Assignor:           c.Assignor,
		Oldest:             c.Oldest,
		EnableTLS:          c.EnableTLS,
		InsecureSkipVerify: c.InsecureSkipVerify,
		ClientCertFile:     c.ClientCertFile,
		ClientKeyFile:      c.ClientKeyFile,
		CACertFile:         c.CACertFile,
	}
}

func GetConfig() *ApplicationConfig {
	cfg := new(ApplicationConfig)
	if err := cfg.loadFromEnv(); err != nil {
		log.Fatal("Failed to load Decision Config", err)
		return nil
	}

	return cfg
}

func (c *ApplicationConfig) loadFromEnv() error {
	return envconfig.Process("", c)
}

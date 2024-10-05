package model

type ConsumerMessage struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
}

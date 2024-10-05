package model

type ConsumerMessage struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
}

type ScoreMessage struct {
	UserId     int64
	Score      int64
	QuizId     int64
	QuestionId int64
}

package service

import "score-service/model"

type ScoreService interface {
	UpdateScore(msg model.ConsumerMessage)
}

type scoreService struct{}

func (s scoreService) UpdateScore(msg model.ConsumerMessage) {
	//TODO implement me
	panic("implement me")
}

func NewScoreService() ScoreService {
	return &scoreService{}
}

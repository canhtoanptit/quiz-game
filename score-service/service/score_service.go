package service

import (
	"context"
	"encoding/json"
	"score-service/model"
	"score-service/proto"
	"score-service/repository"

	log "github.com/sirupsen/logrus"
)

type ScoreService interface {
	UpdateScore(msg model.ConsumerMessage)
	GetQuizScore(ctx context.Context, scoreRequests []*proto.ScoreRequest) ([]*proto.Score, error)
}

type scoreService struct {
	scoreRepository repository.ScoreRepository
}

func (s scoreService) GetQuizScore(ctx context.Context, scoreRequests []*proto.ScoreRequest) ([]*proto.Score, error) {
	return s.scoreRepository.GetScore(ctx, scoreRequests)
}

func (s scoreService) UpdateScore(msg model.ConsumerMessage) {
	var scoreMsg model.ScoreMessage
	err := json.Unmarshal(msg.Value, &scoreMsg)
	if err != nil {
		log.Errorf("[ScoreService] Unmarshal score msg fail of key %s, %v", msg.Key, err)
		return
	}

	err = s.scoreRepository.UpdateScore(context.TODO(), scoreMsg)
	if err != nil {
		log.Errorf("[ScoreService] Update score fail of key %s, %v", msg.Key, err)
	}
}

func NewScoreService(scoreRepo repository.ScoreRepository) ScoreService {
	return &scoreService{
		scoreRepository: scoreRepo,
	}
}

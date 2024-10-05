package repository

import (
	"context"
	"score-service/model"
	"score-service/proto"

	log "github.com/sirupsen/logrus"
)

type ScoreRepository interface {
	UpdateScore(ctx context.Context, score model.ScoreMessage) error
	GetScore(ctx context.Context, scoreRequests []*proto.ScoreRequest) ([]*proto.Score, error)
}

func NewScoreRepository(redisUrl string, mysqlUrl string) ScoreRepository {
	redisRepo, err := NewRedisScoreRepository(redisUrl)
	if err != nil {
		panic(err)
	}

	mysqlClient := NewMysqlClient(mysqlUrl)
	return &scoreRepository{
		redisRepository: redisRepo,
		mysqlClient:     mysqlClient,
	}
}

type scoreRepository struct {
	redisRepository RedisScoreRepository
	mysqlClient     *MysqlClient
}

func (s scoreRepository) GetScore(ctx context.Context, scoreRequests []*proto.ScoreRequest) ([]*proto.Score, error) {
	var result []*proto.Score
	for _, scoreRequest := range scoreRequests {
		score, err := s.redisRepository.GetScore(ctx, scoreRequest.QuizId, scoreRequest.UserId)
		if err != nil {
			return result, err
		}

		result = append(result, &proto.Score{
			QuizId: scoreRequest.QuizId,
			UserId: scoreRequest.UserId,
			Score:  score,
		})
	}

	return result, nil
}

func (s scoreRepository) UpdateScore(ctx context.Context, answer model.ScoreMessage) error {
	err := s.redisRepository.UpdateScore(ctx, answer)
	if err != nil {
		log.Errorf("[ScoreRepository] UpdateScore to redis of msg %+v err: %+v", answer, err)
		// TODO recalculate to update score
	}

	stmt, err := s.mysqlClient.db.Prepare("INSERT INTO quiz_results (quiz_id, user_id, question_id, score) " +
		"VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Errorf("[ScoreRepository] Prepare new score answer to db of msg %+v err: %+v", answer, err)
		return err
	}
	_, err = stmt.Exec(answer.QuizId, answer.UserId, answer.QuestionId, answer.Score)
	if err != nil {
		log.Errorf("[ScoreRepository] Insert new score answer to db of msg %+v err: %+v", answer, err)
		return err
	}

	return nil
}

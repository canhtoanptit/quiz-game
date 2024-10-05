package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"score-service/model"
)

const maxExp = 60 * 60 * 24

type RedisScoreRepository interface {
	UpdateScore(ctx context.Context, answer model.ScoreMessage) error
	GetScore(ctx context.Context, quizId int64, userId int64) (int64, error)
	RemoveScore(ctx context.Context, quizId int64, userId int64) error
}

func NewRedisScoreRepository(redisURL string) (RedisScoreRepository, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	pong := client.Ping(ctx)
	_, err := pong.Result()
	if err != nil {
		return nil, err
	}

	return &redisScoreRepository{
		client,
	}, nil
}

type redisScoreRepository struct {
	redisClient *redis.Client
}

func (r redisScoreRepository) UpdateScore(ctx context.Context, answer model.ScoreMessage) error {
	answerKey := getAnswerKeyUnique(answer)
	hasAnswer, err := r.redisClient.Exists(ctx, answerKey).Result()
	if err != nil {
		return err
	}
	if hasAnswer > 0 {
		return nil
	}
	key := getRedisKeyScore(answer.QuizId, answer.UserId)
	return r.redisClient.IncrBy(ctx, key, answer.Score).Err()
}

func (r redisScoreRepository) GetScore(ctx context.Context, quizId int64, userId int64) (int64, error) {
	key := getRedisKeyScore(quizId, userId)
	return r.redisClient.Get(ctx, key).Int64()
}

func (r redisScoreRepository) RemoveScore(ctx context.Context, quizId int64, userId int64) error {
	key := getRedisKeyScore(quizId, userId)
	return r.redisClient.Del(ctx, key).Err()
}

func getAnswerKeyUnique(answer model.ScoreMessage) string {
	return fmt.Sprintf("quizanswer:%d_%d_%d", answer.QuizId, answer.UserId, answer.QuestionId)
}

func getRedisKeyScore(quizId int64, userId int64) string {
	return fmt.Sprintf("quizscore:%d_%d", quizId, userId)
}

package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"score-service/proto"
	"score-service/service"
)

// GrpcServer Server implements the QuizService gRPC service
type GrpcServer struct {
	proto.UnimplementedQuizServiceServer
	scoreService service.ScoreService
}

func NewGrpcServer(scoreService service.ScoreService) *GrpcServer {
	return &GrpcServer{
		scoreService: scoreService,
	}
}

// GetScore retrieves the scores of all users for a specific quiz
func (s *GrpcServer) GetScore(ctx context.Context, req *proto.GetScoreRequest) (*proto.GetScoreResponse, error) {
	quizScores, err := s.scoreService.GetQuizScore(ctx, req.ScoreRequests)

	if err != nil {
		log.Errorf("[GrpcServer] GetScore req %+v failed: %v", req, err)
		return nil, err
	}
	return &proto.GetScoreResponse{
		Scores: quizScores,
	}, nil
}

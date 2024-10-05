package server

import (
	"context"
	"score-service/proto"
)

// GrpcServer Server implements the QuizService gRPC service
type GrpcServer struct {
	proto.UnimplementedQuizServiceServer
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

// GetScore retrieves the scores of all users for a specific quiz
func (s *GrpcServer) GetScore(ctx context.Context, req *proto.GetScoreRequest) (*proto.GetScoreResponse, error) {
	var scores []*proto.Score

	return &proto.GetScoreResponse{
		Scores: scores,
	}, nil
}

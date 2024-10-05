package main

import (
	"fmt"
	"net"
	"score-service/config"
	"score-service/pkg/kafka"
	"score-service/pkg/server"
	"score-service/proto"
	"score-service/repository"
	"score-service/service"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	log.Info("Event-Connector Application starting...")
	err := godotenv.Load("config.env")
	if err != nil {
		panic("Error loading .env file")
	}

	appConf := config.GetConfig()

	userSubmitCG := kafka.InitSaramaConsumerAndPanicIfError(appConf.KafkaConfig.ToConfig(
		appConf.UserSubmitAnswerTopic, appConf.ScoreServiceConsumerGroup))

	scoreRepo := repository.NewScoreRepository(appConf.RedisUrl, appConf.MysqlUrl)
	scoreService := service.NewScoreService(scoreRepo)
	go func() {
		errHandleUpdateScore := userSubmitCG.Consume(scoreService.UpdateScore)
		if errHandleUpdateScore != nil {
			panic(errHandleUpdateScore)
		}
	}()

	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the quiz service with the server
	proto.RegisterQuizServiceServer(grpcServer, server.NewGrpcServer(scoreService))

	// Start the gRPC server
	fmt.Println("gRPC server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

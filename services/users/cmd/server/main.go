package main

import (
	"fmt"
	"net"

	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
	userConfig "github.com/0xsj/fn-go/services/users/internal/config"
	"github.com/0xsj/fn-go/services/users/internal/handler"
	"google.golang.org/grpc"
)

func main() {
	envProvider := config.NewEnvProvider("USERS")
	cfg := userConfig.LoadConfig(envProvider)

	baseLogger := logging.NewSimpleLogger(logging.InfoLevel)
	logger := baseLogger.With(logging.F("service", "users"))
	logger.Info("Starting users service")

	server := grpc.NewServer()

	userHandler := handler.NewUserServiceHandler(cfg, logger)
	users.RegisterUserServiceServer(server, userHandler)

	port := cfg.Server.Port
	if port == 0 {
		port = 50051
	}

	logger.Info("Users service listening", logging.F("port", port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("Failed to listen", logging.F("error", err))
	}

	if err := server.Serve(lis); err !=  nil {
		logger.Fatal("Failed to serve", logging.F("error", err))
	}
}
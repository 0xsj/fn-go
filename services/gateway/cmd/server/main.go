package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	userClient users.UserServiceClient
	logger	logging.Logger
}

func main() {
    port := 8000
    if portEnv := os.Getenv("GATEWAY_PORT"); portEnv != "" {
        if p, err := strconv.Atoi(portEnv); err == nil {
            port = p
        }
    }

    usersServiceAddress := "localhost:50051"
    if addr := os.Getenv("USERS_SERVICE_ADDRESS"); addr != "" {
        usersServiceAddress = addr
    }

    baseLogger := logging.NewSimpleLogger(logging.InfoLevel)
    serviceLogger := baseLogger.With(logging.F("service", "gateway"))
    serviceLogger.Info("Starting API gateway")

    serviceLogger.Info("Connecting to users service", logging.F("address", usersServiceAddress))
    userConn, err := grpc.Dial(usersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        serviceLogger.Fatal("Failed to connect to users service", logging.F("error", err))
    }
    defer userConn.Close()

    userClient := users.NewUserServiceClient(userConn)

    // Create gateway
    gateway := &Gateway{
        userClient: userClient,
        logger:     serviceLogger,
    }

    // Set up HTTP handlers
    http.HandleFunc("/users", gateway.handleUsers)
    http.HandleFunc("/users/", gateway.handleUserByID)

    // Start the server
    addr := fmt.Sprintf(":%d", port)
    serviceLogger.Info("Gateway listening", logging.F("address", addr))
    if err := http.ListenAndServe(addr, nil); err != nil {
        serviceLogger.Fatal("Failed to start HTTP server", logging.F("error", err))
    }
}
// services/gateway/cmd/server/main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
	"github.com/0xsj/fn-go/services/gateway/internal/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

    userConn, err := grpc.NewClient(
        usersServiceAddress,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        serviceLogger.Fatal("Failed to connect to users service", logging.F("error", err))
    }
    defer userConn.Close()

    userClient := users.NewUserServiceClient(userConn)

    userHandler := handler.NewUserHandler(userClient, serviceLogger)

    http.HandleFunc("/users", userHandler.HandleUsers)
    http.HandleFunc("/users/", userHandler.HandleUserByID)
    addr := fmt.Sprintf(":%d", port)
    serviceLogger.Info("Gateway listening", logging.F("address", addr))
    if err := http.ListenAndServe(addr, nil); err != nil {
        serviceLogger.Fatal("Failed to start HTTP server", logging.F("error", err))
    }
}
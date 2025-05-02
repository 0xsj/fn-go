package handler

import (
	"context"

	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
	"github.com/0xsj/fn-go/services/users/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceHandler struct {
	users.UnimplementedUserServiceServer
	cfg *config.Config
	logger logging.Logger
}

func NewUserServiceHandler(cfg *config.Config, logger logging.Logger) *UserServiceHandler {
	return &UserServiceHandler{
		cfg: cfg,
		logger: logger,
	}
}

func (h *UserServiceHandler) GetUser(ctx context.Context, req *users.GetUserRequest) (*users.GetUserResponse, error) {
	h.logger.Info("GetUser request received", logging.F("id", req.Id))
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	user := &users.User{
		Id: req.Id,
		Email: "user@example.com",
		Name: "Example User",
		Role: "user",
	}

	return &users.GetUserResponse{
		User: user,
	}, nil
}
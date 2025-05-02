package handler

import (
	"log"

	"github.com/0xsj/fn-go/pkg/proto/users"
	"github.com/0xsj/fn-go/services/users/internal/config"
)

type UserServiceHandler struct {
	users.UnimplementedUserServiceServer
	cfg *config.Config
	logger *log.Logger
}

func NewServiceHandler(cfg *config.Config, logger *log.Logger) *UserServiceHandler {
	return &UserServiceHandler{
		cfg: cfg,
		logger: logger,
	}
}


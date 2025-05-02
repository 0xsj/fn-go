package handler

import (
	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
	"github.com/0xsj/fn-go/services/users/internal/config"
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


package handler

import (
	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
)

type UserHandler struct {
	client users.UserServiceClient
	logger logging.Logger
}

func NewUserHandler(client users.UserServiceClient, logger logging.Logger) *UserHandler {
	return &UserHandler{
		client: client,
		logger: logger.With(logging.F("component", "user-handler")),
	}
}
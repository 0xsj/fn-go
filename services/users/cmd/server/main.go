package main

import (
	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/logging"
	userConfig "github.com/0xsj/fn-go/services/users/internal/config"
)

func main() {
	envProvider := config.NewEnvProvider("USERS")
	cfg := userConfig.LoadConfig(envProvider)

	logger := logging.NewSimpleLogger(logging.InfoLevel)


}
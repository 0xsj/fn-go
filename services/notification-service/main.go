package main

import (
	"github.com/0xsj/fn-go/pkg/common/log"
)


func main(){
	logger := log.Default()
	logger = logger.WithLayer("notification-service")
	logger.Info("Initializing notification service")
}
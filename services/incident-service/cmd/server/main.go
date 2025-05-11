package main

import (
	"github.com/0xsj/fn-go/pkg/common/log"
)


func main(){
	logger := log.Default()
	logger = logger.WithLayer("incident-service")
	logger.Info("Initializing incident service")
}
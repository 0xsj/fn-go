package main

import (
	"github.com/0xsj/fn-go/pkg/common/log"
)

func main(){
	logger := log.Default()
	logger = logger.WithLayer("entity-service")
	logger.Info("Initializing entity service")
}
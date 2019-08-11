package main

import (
	"github.com/Russiancold/testApp/api"
	"github.com/Russiancold/testApp/service"
	"gopkg.in/tylerb/graceful.v1"
	"time"
)

func main() {
	service.Init()
	defer service.GetService().Close()
	graceful.Run(":8080", time.Second * 10, api.New())
}

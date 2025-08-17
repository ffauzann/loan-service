package main

import (
	"github.com/ffauzann/loan-service/internal/app"
)

var cfg app.Config

func init() {
	cfg.Setup()
}

func main() {
	cfg.StartServer()
}

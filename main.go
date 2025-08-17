package main

import (
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/app"
)

var cfg app.Config

func init() {
	cfg.Setup()
}

func main() {
	cfg.StartServer()
}

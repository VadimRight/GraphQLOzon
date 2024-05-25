package main

import (
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"github.com/VadimRight/GraphQLOzon/api"
	"github.com/VadimRight/GraphQLOzon/storage"
)

func main() {
	cfg := bootstrap.LoadConfig()
	storage := storage.InitPostgresDatabase(cfg)
	defer storage.ClosePostgres(storage)
	api.InitServer(cfg, storage)
}

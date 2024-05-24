package main

import (
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"github.com/VadimRight/GraphQLOzon/api"
)

func main() {
	cfg := bootstrap.LoadConfig()
	storage := bootstrap.InitPostgresDatabase(cfg)
	defer bootstrap.CloseDB(storage)
	api.InitServer(cfg, storage)
}

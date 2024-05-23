package main

import (
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"github.com/VadimRight/GraphQLOzon/api"
)

func main() {
	cfg := bootstrap.LoadConfig()
	db := bootstrap.InitPostgresDatabase(cfg)
	defer bootstrap.CloseDB(db)
	api.InitServer(cfg)	
}

package main

import (
	"github.com/VadimRight/GraphQLOzon/api"
	"github.com/VadimRight/GraphQLOzon/internal/config"
	"github.com/VadimRight/GraphQLOzon/storage"
)

func main() {
	cfg := config.LoadConfig()
	storageType := storage.StorageType(cfg)
	defer func() {
		if postgresStorage, ok := storageType.(*storage.PostgresStorage); ok {
			postgresStorage.ClosePostgres()
		}
	}()
	api.InitServer(cfg, storageType)
}

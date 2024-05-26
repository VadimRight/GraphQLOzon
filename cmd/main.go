package main

import (
	"github.com/VadimRight/GraphQLOzon/api"
	"github.com/VadimRight/GraphQLOzon/internal/config"
	"github.com/VadimRight/GraphQLOzon/storage"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация хранилища данных
	dbStorage := storage.InitPostgresDatabase(cfg)
	defer dbStorage.ClosePostgres()

	// Инициализация и запуск сервера
	api.InitServer(cfg, dbStorage)
}

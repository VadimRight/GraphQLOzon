package main

import (
	"github.com/VadimRight/OzonTestTask/bootstrap"
)

func main() {
	cfg := bootstrap.LoadConfig()
	db := bootstrap.InitPostgresDatabase(cfg)
	defer bootstrap.CloseDB(db)
}

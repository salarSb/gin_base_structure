package main

import (
	"base_structure/src/api"
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/data/cache"
	"base_structure/src/data/db"
	"base_structure/src/data/db/migrations"
)

func main() {
	cfg := config.GetConfig()
	constants.InitConstants()
	cache.GetRedis(cfg)
	defer cache.CloseRedis()
	db.GetDb(cfg)
	migrations.Up1(cfg)
	defer db.CloseDb()
	api.InitServer(cfg)
}

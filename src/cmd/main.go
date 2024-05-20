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
	cache.GetRedis()
	defer cache.CloseRedis()
	db.GetDb()
	migrations.Up1()
	defer db.CloseDb()
	api.InitServer(cfg)
}

package main

import (
	"base_structure/src/api"
	"base_structure/src/config"
	"base_structure/src/data/cache"
	"base_structure/src/data/db"
)

func main() {
	cfg := config.GetConfig()
	//logger := logging.NewLogger(cfg)
	err := cache.InitRedis(cfg)
	if err != nil {
		//logger.Fatal(logging.Redis, logging.StartUp, err.Error(), nil)
		return
	}
	defer cache.CloseRedis()
	err = db.InitDb(cfg)
	if err != nil {
		//logger.Fatal(logging.Postgres, logging.StartUp, err.Error(), nil)
		return
	}
	//migrations.Up1()
	defer db.CloseDb()
	api.InitServer(cfg)
}

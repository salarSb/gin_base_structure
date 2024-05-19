package main

import (
	"base_structure/src/api"
	"base_structure/src/config"
	"base_structure/src/data/cache"
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
	api.InitServer(cfg)
}

package main

import (
	"base_structure/src/api"
	"base_structure/src/config"
)

func main() {
	cfg := config.GetConfig()
	api.InitServer(cfg)
}

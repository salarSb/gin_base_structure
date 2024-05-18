package api

import (
	"base_structure/src/api/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	RegisterRoutes(r)
	//err := r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))
	err := r.Run(fmt.Sprintf(":5005"))
	if err != nil {
		//logger.Fatal(logging.Internal, logging.Api, "error on running router", nil)
		return
	}
}

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		health := v1.Group("/health")
		routers.Health(health)
	}
}

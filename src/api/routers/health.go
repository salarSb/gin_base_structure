package routers

import (
	"base_structure/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	h := handlers.NewHealthHandler()
	r.GET("/", h.Health)
}

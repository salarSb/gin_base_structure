package routers

import (
	"base_structure/src/api/handlers"
	"base_structure/src/api/middlewares"
	"base_structure/src/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewUserHandler(cfg)

	// public endpoints
	router.POST("/send-otp", middlewares.OtpLimiter(cfg), h.SendOtp)
	router.POST("/login-by-username", h.LoginByUsername)
	router.POST("/register-by-username", h.RegisterByUsername)
	router.POST("/login-by-mobile", h.RegisterLoginByMobileNumber)

	// protected endpoints
	auth := router.Group("").Use(middlewares.Authentication(cfg))
	auth.POST("logout", h.Logout)
}

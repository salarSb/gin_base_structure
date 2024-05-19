package api

import (
	"base_structure/src/api/middlewares"
	"base_structure/src/api/routers"
	"base_structure/src/api/validations"
	"base_structure/src/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer(cfg *config.Config) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.Cors(cfg))
	RegisterRoutes(r)
	RegisterValidators()
	err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
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

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validations.IranianMobileNumberValidator, true)
		if err != nil {
			//logger.Fatal(
			//	logging.Validation,
			//	logging.MobileValidation,
			//	"Error on registering custom mobile validation",
			//	nil,
			//)
			return
		}
		err = val.RegisterValidation("password", validations.PasswordValidator, true)
		if err != nil {
			//logger.Fatal(
			//	logging.Validation,
			//	logging.PasswordValidation,
			//	"Error on registering custom password validation",
			//	nil,
			//)
			return
		}
	}
}

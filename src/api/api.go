package api

import (
	"base_structure/docs"
	"base_structure/src/api/middlewares"
	"base_structure/src/api/routers"
	"base_structure/src/api/validations"
	"base_structure/src/config"
	"base_structure/src/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

func InitServer(cfg *config.Config) {
	logger := logging.NewLogger(cfg)
	defer func(logger logging.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Fatal(logging.Internal, logging.Closing, "error on syncing logger", nil)
		}
	}(logger)
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(logging.Internal, logging.Api, "error on reading .env", nil)
		return
	}
	appEnv := os.Getenv("APP_ENV")
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()
	if appEnv == "development" {
		r.Use(gin.Logger(), gin.CustomRecovery(middlewares.ErrorHandler))
	} else {
		r.Use(gin.Logger(), gin.Recovery())
	}
	r.Use(middlewares.Cors(cfg))
	r.Use(middlewares.StructuredLogger(logger))
	RegisterValidators(logger)
	RegisterRoutes(r, cfg)
	RegisterSwagger(r, cfg)
	err = r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		logger.Fatal(logging.Internal, logging.Api, "error on running router", nil)
		return
	}
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		//User
		users := v1.Group("/users")

		//User
		routers.User(users, cfg)
	}
}

func RegisterValidators(logger logging.Logger) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("ir_mobile", validations.IranianMobileNumberValidator, true)
		if err != nil {
			logger.Fatal(
				logging.Validation,
				logging.MobileValidation,
				"Error on registering custom mobile validation",
				nil,
			)
			return
		}
		err = val.RegisterValidation("password", validations.PasswordValidator, true)
		if err != nil {
			logger.Fatal(
				logging.Validation,
				logging.PasswordValidation,
				"Error on registering custom password validation",
				nil,
			)
			return
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "base structure web api golang-gin"
	docs.SwaggerInfo.Description = "use this base structure to expand your projects in golang using gin framework and gorm orm. the base structure uses postgres as database and has authentication and authorization service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

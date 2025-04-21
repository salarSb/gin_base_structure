package db

import (
	"base_structure/src/config"
	"base_structure/src/pkg/logging"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

var dbClient *gorm.DB
var dbInit sync.Once
var logger = logging.NewLogger(config.GetConfig())

func InitDb(cfg *config.Config) error {
	var err error
	dbInit.Do(func() {
		cnn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.DbName,
			cfg.Postgres.SSLMode,
		)
		dbClient, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
		if err != nil {
			return
		}
		sqlDb, _ := dbClient.DB()
		err = sqlDb.Ping()
		if err != nil {
			return
		}
		sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)
		sqlDb.SetMaxOpenConns(cfg.Postgres.MaxOpenConnections)
		sqlDb.SetConnMaxLifetime(cfg.Postgres.ConnectionMaxLifetime * time.Minute)
		logger.Info(logging.Postgres, logging.StartUp, "db connection established", nil)
	})
	return err
}

func GetDb(cfg *config.Config) *gorm.DB {
	if dbClient == nil {
		err := InitDb(cfg)
		if err != nil {
			logger.Fatal(logging.Postgres, logging.StartUp, err.Error(), nil)
		}
	}
	return dbClient
}

func CloseDb() {
	if dbClient != nil {
		cnn, _ := dbClient.DB()
		err := cnn.Close()
		if err != nil {
			logger.Info(logging.Postgres, logging.Closing, "error on closing db connection", nil)
			return
		}
		dbClient = nil
	}
}

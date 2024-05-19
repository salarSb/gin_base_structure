package db

import (
	"base_structure/src/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var dbClient *gorm.DB

//var logger = logging.NewLogger(config.GetConfig())

func InitDb(cfg *config.Config) error {
	var err error
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
		return err
	}
	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxOpenConnections)
	sqlDb.SetConnMaxLifetime(cfg.Postgres.ConnectionMaxLifetime * time.Minute)
	//logger.Info(logging.Postgres, logging.StartUp, "db connection established", nil)
	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	cnn, _ := dbClient.DB()
	err := cnn.Close()
	if err != nil {
		//logger.Info(logging.Postgres, logging.Closing, "error on closing db connection", nil)
		return
	}
}

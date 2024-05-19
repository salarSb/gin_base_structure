package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Password PasswordConfig
	Cors     CorsConfig
	Otp      OtpConfig
	Jwt      JwtConfig
}

type ServerConfig struct {
	Port    string
	RunMode string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

type PostgresConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	DbName                string
	SSLMode               string
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	Db                 int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleCheckFrequency time.Duration
	IdleTimeout        time.Duration
}

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type CorsConfig struct {
	AllowOrigins string
}

type OtpConfig struct {
	ExpireTime time.Duration
	Digits     int
	Limiter    time.Duration
}

type JwtConfig struct {
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
	Secret                     string
	RefreshSecret              string
}

func getConfigPath(env string) (string, error) {
	if env == "docker" {
		return "/app/config/config-docker", nil
	} else if env == "production" {
		return "/config/config-production", nil
	} else if env == "development" {
		return "config/config-development", nil
	}
	return "", errors.New("wrong environment")
}

func LoadConfig(filename string, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath("./src")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		log.Printf(" Unable to read config: %v", err)
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf(" Unable to parse config: %v", err)
		return nil, err
	}
	return &cfg, nil
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error in load .env file, %v", err)
	}
	cfgPath, err := getConfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatalf("Error in finding config path, %v", err)
	}
	v, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatalf("Error in load config, %v", err)
	}
	cfg, err := ParseConfig(v)
	if err != nil {
		log.Fatalf("Error in parse config, %v", err)
	}
	return cfg
}

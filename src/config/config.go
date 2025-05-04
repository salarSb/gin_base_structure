package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	cfg  *Config
	once sync.Once
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

func LoadDotEnv() {
	startDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working dir: %v", err)
	}
	wd := startDir
	for {
		envPath := filepath.Join(wd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if lerr := godotenv.Load(envPath); lerr != nil {
				log.Fatalf("unable to load %s: %v", envPath, lerr)
			}
			return
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			log.Fatalf(".env not found (searched upward from %s)", startDir)
		}
		wd = parent
	}
}

func GetConfig() *Config {
	once.Do(func() {
		LoadDotEnv()
		v, err := resolveConfig()
		if err != nil {
			log.Fatalf("config error: %v", err)
		}
		c, err := ParseConfig(v)
		if err != nil {
			log.Fatalf("parsing config: %v", err)
		}
		cfg = c
	})
	return cfg
}

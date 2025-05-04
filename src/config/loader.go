package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var cfgFlag = flag.String("config", "", "absolute path to YAML/TOML config")

func resolveConfig() (*viper.Viper, error) {
	if *cfgFlag != "" {
		return loadAbsolute(*cfgFlag)
	}
	if cf := os.Getenv("CONFIG_FILE"); cf != "" {
		return loadAbsolute(cf)
	}
	return nil, fmt.Errorf("no config specified: use --config or set CONFIG_FILE")
}

func loadAbsolute(full string) (*viper.Viper, error) {
	ext := strings.TrimPrefix(filepath.Ext(full), ".")
	name := strings.TrimSuffix(full, "."+ext)
	return LoadConfig(name, ext)
}

func LoadConfig(relPath, fileType string) (*viper.Viper, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot get cwd: %w", err)
	}
	for {
		full := filepath.Join(wd, fmt.Sprintf("%s.%s", relPath, fileType))
		if _, err := os.Stat(full); err == nil {
			v := viper.New()
			v.SetConfigFile(full)
			v.AutomaticEnv()

			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("unable to read config: %w", err)
			}
			return v, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return nil, fmt.Errorf("config file %s.%s not found", relPath, fileType)
		}
		wd = parent
	}
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

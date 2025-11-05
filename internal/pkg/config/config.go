package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	httpConfig "github.com/yrss1/workout/internal/pkg/http/config"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" toml:"log_level" env-default:"debug"`
	HTTP     httpConfig.HTTPConfig
}

func ParseConfig(configPath string) (*Config, error) {
	config := &Config{}

	var err error

	if configPath != "" {
		err = cleanenv.ReadConfig(configPath, config)
	} else {
		_ = godotenv.Load()
		err = cleanenv.ReadEnv(config)
	}
	fmt.Println(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

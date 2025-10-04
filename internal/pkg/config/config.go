package config

import (
	"github.com/ilyakaznacheev/cleanenv"
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
		err = cleanenv.ReadEnv(config)
	}

	if err != nil {
		return nil, err
	}
	return config, nil
}

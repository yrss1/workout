package container

import (
	"github.com/yrss1/workout/internal/pkg/config"
	"github.com/yrss1/workout/internal/pkg/log"
)

type App struct {
	config *config.Config
	logger *log.Log
}

func NewApp(configPath string) (*App, error) {
	config, err := config.ParseConfig(configPath)
	if err != nil {
		return nil, err
	}
	logger, err := log.NewLog(config.LogLevel)
	if err != nil {
		return nil, err
	}

	var container = &App{
		config: config,
		logger: logger,
	}

	return container, nil
}

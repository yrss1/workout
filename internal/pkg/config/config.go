package config

import (
	httpConfig "workout/internal/pkg/http/config"
)

type Config struct {
	AppName string
	HTTP    httpConfig.Config
}

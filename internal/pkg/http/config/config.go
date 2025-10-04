package config

type HTTPConfig struct {
	BindAddr         string   `env:"BIND_ADDR" toml:"bind_addr" env-default:":8000"`
	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" toml:"allowed_origins"`
	AllowedMethods   []string `env:"ALLOWED_METHODS" toml:"allowed_methods"`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" toml:"allowed_headers"`
	AllowCredentials bool     `env:"ALLOW_CREDENTIALS" toml:"allow_credentials"`
}

package config

import "os"

type Config struct {
	AppEnv    string
	JWTSecret string
}

func Load() Config {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" {
		appEnv = "development"
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if appEnv != "development" && jwtSecret == "" {
		panic("JWT_SECRET is required outside development environment")
	}

	// Secret only used for local development.
	if jwtSecret == "" {
		jwtSecret = "development-secret-only"
	}

	return Config{
		AppEnv:    appEnv,
		JWTSecret: jwtSecret,
	}
}

func (c Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct { // структура конфио котоырй содержит в себе данные с env  чтобы использовать в проекте
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	ApiServerPort    string `env:"API_SERVER_PORT"`
	ApiServerHost    string `env:"API_SERVER_HOST"`
	DatabasePortTest string `env:"DB_PORT_TEST"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePass     string `env:"DB_PASSWORD"`
	Environment      Env    `env:"ENV" envDefault:"dev"`
	ProjectRoot      string `env:"PROJECT_ROOT"`
	JwtSecret        string `env:"JWT_SECRET"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

type Env string

const (
	EnvDev  Env = "dev"
	EnvTest Env = "test"
)

func (c *Config) DatabaseURL() string {
	port := c.DatabasePort
	if c.Environment == EnvTest {
		port = c.DatabasePortTest
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DatabaseUser, c.DatabasePass, c.DatabaseHost, port, c.DatabaseName)
}

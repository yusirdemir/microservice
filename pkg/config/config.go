package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const DefaultPort = "3000"

type Config struct {
	App      AppConfig      `yaml:"app" env-prefix:"APP_"`
	Server   ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
	Logger   LoggerConfig   `yaml:"logger" env-prefix:"LOGGER_"`
	Database DatabaseConfig `yaml:"database" env-prefix:"DATABASE_"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver" env:"DRIVER" env-default:"memory"`
	Host     string `yaml:"host" env:"HOST"`
	Bucket   string `yaml:"bucket" env:"BUCKET"`
	Username string `yaml:"username" env:"USERNAME"`
	Password string `yaml:"password" env:"PASSWORD"`
}

type ServerConfig struct {
	ReadTimeout  string `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout string `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  string `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"120s"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"NAME" env-default:"User Service"`
	Port string `yaml:"port" env:"PORT" env-default:"3000"`
	Env  string `yaml:"env" env:"ENV" env-default:"development"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LEVEL" env-default:"info"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Try to look for config files in logical places
	paths := []string{
		fmt.Sprintf("config.%s.yaml", env),
		fmt.Sprintf("config/config.%s.yaml", env),
	}

	var configPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			configPath = p
			break
		}
	}

	// Load from file if exists, otherwise load from env
	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to read config file '%s': %w", configPath, err)
		}
	} else {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("failed to read env: %w", err)
		}
	}

	return cfg, nil
}

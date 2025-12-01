package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const DefaultPort = "3000"

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
}

type ServerConfig struct {
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	IdleTimeout  string `mapstructure:"idle_timeout"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")

	v.SetDefault("app.port", DefaultPort)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	cfg.App.Env = env

	return &cfg, nil
}

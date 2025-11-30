package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Logger LoggerConfig `mapstructure:"logger"`
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
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config not found")
		}
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}

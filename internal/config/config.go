package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DSN                     string        `mapstructure:"DSN"`
	PORT                    string        `mapstructure:"PORT"`
	TOKEN_SYMMETRIC_KEY     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

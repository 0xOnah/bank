package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DSN                     string        `mapstructure:"DSN"`
	HTTP_SERVER_ADDRESS     string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPC_SERVER_ADDRESS     string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TOKEN_SYMMETRIC_KEY     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATATION"`
	REFRESH_TOKEN_DURATION  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	ENVIRONMENT             string        `mapstructure:"ENVIRONMENT"`
	LOG_LEVEL               string        `mapstructure:"LOG_LEVEL"`
	REDIS_ADDRESS           string        `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	//reading from enviroment varaibles
	if err = viper.BindEnv("DSN"); err != nil {
		return Config{}, fmt.Errorf("bind 'DSN' env variable: %w", err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("no config file found or error reading config: %v", err)
	}

	err = viper.Unmarshal(&config)
	return
}

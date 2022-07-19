package config

import (
	"errors"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Login string `mapstructure:"BITLY_OAUTH_LOGIN"`
	Token string `mapstructure:"BITLY_OAUTH_TOKEN"`
	Port  string `mapstructure:"PORT"`
}

var config = Config{}

func isBlankAny(cfg Config) bool {
	v := reflect.ValueOf(cfg)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			return true
		}
	}
	return false
}

func Read(path string) error {

	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	// Don't throw error cause variables still might have been set manually without .env
	viper.ReadInConfig()

	config = Config{
		Login: viper.GetString("BITLY_OAUTH_LOGIN"),
		Token: viper.GetString("BITLY_OAUTH_TOKEN"),
		Port:  viper.GetString("PORT"),
	}

	if isBlankAny(config) {
		return errors.New("some variables are not initialized")
	}

	return nil
}

func Get() *Config {
	return &config
}

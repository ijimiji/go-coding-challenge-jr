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
var read = false

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
	if read {
		return nil
	}

	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil || isBlankAny(config) {
		return errors.New("can't read config correctly")
	}
	read = true
	return nil
}

func Get() *Config {
	return &config
}

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func MustLoad() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("CONFIG_PATH")) // /home/cat/test_area/go/anonymous_chat

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

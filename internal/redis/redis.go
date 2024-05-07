package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Client *redis.Client

func MustLoad() {
	Client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}

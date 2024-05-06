package redis_client

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func MustLoad() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}

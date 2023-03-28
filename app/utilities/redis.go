package utilities

import (
	"github.com/go-redis/redis/v9"
)

var Redis redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: "vir-pc-redis:6379",
		// TODO get pass from db
		Password: "password123",
		DB:       0,
	})
	Redis = *client
}

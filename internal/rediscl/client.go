package rediscl

import (
    "fmt"
    "github.com/redis/go-redis/v9"
)

type RedisClient struct {
    Client *redis.Client
}

var RDB *RedisClient

func InitRedisClient(host, port string) {
    RDB =  &RedisClient{
            Client: redis.NewClient(&redis.Options{
                Addr: fmt.Sprintf("%s:%s", host, port),
                Password: "",
                DB: 0,
        }),
    }
}




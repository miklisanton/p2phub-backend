package rediscl

import (
    "fmt"
    "github.com/redis/go-redis/v9"
    "context"
)

type RedisClient struct {
    Client *redis.Client
    Ctx    context.Context
}

var RDB *RedisClient

func InitRedisClient(host, port string) {
    RDB =  &RedisClient{
            Ctx: context.Background(),
            Client: redis.NewClient(&redis.Options{
                Addr: fmt.Sprintf("%s:%s", host, port),
                Password: "",
                DB: 0,
        }),
    }
}





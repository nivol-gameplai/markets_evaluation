package configs

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func InitializeRedis() (*redis.Pool, error) {
	redisHost := "10.60.48.4"
	redisPort := "6379"
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 20
	return &redis.Pool{
		MaxIdle: maxConnections,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, fmt.Errorf("redis.Dial: %v", err)
			}
			return c, err
		},
	}, nil
}

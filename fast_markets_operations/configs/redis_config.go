package configs

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
)

/*
// ==================================================================
// https://firebase.google.com/docs/admin/setup
// ==================================================================
// CreateClient /*
//initialization of the redis client (memorystore)
//by reading host and port as environment variable
//INPUTS:
//	redis Pool
//OUTPUTS:
//	initialized redis pool
*/
func InitializeRedis() (*redis.Pool, error) {
	redisHost := os.Getenv("REDISHOST")
	//redisHost := "10.60.48.4"
	//redisPort := "6379"
	redisPort := os.Getenv("REDISPORT")
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

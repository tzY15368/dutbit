package RedisConnector

import (
	"github.com/garyburd/redigo/redis"
	"github.com/wonderivan/logger"
)

func RedisInit() redis.Conn {
	cli, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("Bit_redis_123"))
	if err != nil {
		panic("redis connection error")
	}
	logger.Info("redis connected")
	return cli
}

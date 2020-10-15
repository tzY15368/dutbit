package main

import (
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var redis_conn redis.Conn

func ErrorHandler(err error, ErrorType string) {
	if err != nil {
		logger.Error(ErrorType)
		panic(err)
	}
}
func RedisInit() {
	cli, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("")) //password here
	ErrorHandler(err, "redis connection error")
	logger.Info("redis connected")
	redis_conn = cli
}
func RedirectHandler(c *gin.Context) {
	id := c.Param("id")
	shortlink := "https://sl.dutbit.com/" + id
	//logger.Info(shortlink)
	result, _ := redis.String(redis_conn.Do("get", shortlink))
	c.Redirect(http.StatusTemporaryRedirect, result)
}
func main() {
	RedisInit()
	g := gin.Default()
	g.GET("/:id", RedirectHandler)
	g.Run("127.0.0.1:8820")
	defer redis_conn.Close()
}

package db

import (
	. "utils"

	"github.com/garyburd/redigo/redis"
	"github.com/wonderivan/logger"
)

var redis_conn redis.Conn

func init() {
	cli, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("Bit_redis_123"))
	ErrorHandler(err, "redis connection error")
	logger.Info("redis connected")
	redis_conn = cli
}
func PingRedis() {
	_, err := redis_conn.Do("PING")
	if err != nil {
		panic(err)
	}
	logger.Info("ping ok")
}
func LoggedIn(sessionid string) bool {
	result, _ := redis.Bool(redis_conn.Do("exists", sessionid))
	return result
}
func IsAdmin(sessionid string) bool {
	site, _ := redis.String(redis_conn.Do("hget", sessionid, "site"))
	site_map := JSONToMap(site)
	_, ok := site_map["super_admin"]
	return ok
}
func FetchAllInfo(sessionid string) map[string]string {
	result, err := redis.StringMap(redis_conn.Do("hgetall", sessionid))
	ErrorHandler(err, "fetch all info error")
	result["password"] = ""
	return result
}
func SessionStartRedis(sessionid string, sessioninfo map[string]string) {

}
func DeleteSessionRedis(sessionid string) {
	result, err := redis.Int(redis_conn.Do("del", sessionid))
	if result != 1 {
		logger.Error(err)
	}
}

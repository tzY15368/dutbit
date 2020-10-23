package handlers

import (
	"db"
	. "models"
	"net/http"
	"time"
	. "utils"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

// func session_start(c *gin.Context, sessioninfo map[string]interface{}) string {
// 	COOKIE_URL := ".dutbit.com"
// 	COOKIE_EXPIRE_TIME := 2592000 //seconds
// 	IpArray, _ := c.Request.Header["X-Real-Ip"]
// 	Ip := IpArray[0]
// 	sessionid := GetSessionId(Ip)
// 	//_, err :=
// 	redis_conn.Do("hmset", redis.Args{}.Add(sessionid).AddFlat(sessioninfo)...)
// 	redis_conn.Do("EXPIRE", sessionid, COOKIE_EXPIRE_TIME)
// 	//ErrorHandler(err, "hmset error(redis)")
// 	c.SetCookie("SESSIONID", sessionid, COOKIE_EXPIRE_TIME, "/", COOKIE_URL, true, false)
// 	return sessionid
// }
func UserInfoGetHandler(c *gin.Context) {
	sessionid, _ := c.Cookie("SESSIONID")
	result := db.FetchAllInfo(sessionid)
	c.JSON(200, result)
}
func LogoutHandler(c *gin.Context) {
	COOKIE_URL := ".dutbit.com"
	LOGOUT_REDIRECT_URL := "https://www.dutbit.com/wp20/index.php?logout"
	sessionid, _ := c.Cookie("SESSIONID")
	go db.DeleteSessionRedis(sessionid)
	go db.DeleteSessionMongo(sessionid) // todo
	logger.Info(sessionid, "logging out")
	c.SetCookie("SESSIONID", "removed", -1, "/", COOKIE_URL, false, true)
	c.Redirect(http.StatusTemporaryRedirect, LOGOUT_REDIRECT_URL)
	return
}
func LoginHandler(c *gin.Context) {
	COOKIE_URL := ".dutbit.com"
	COOKIE_EXPIRE_TIME := 2592000 //seconds
	var JSONInput LoginRequest
	ip, _ := c.Request.Header["X-Real-Ip"]
	sessionid := GetSessionId(ip[0])
	now := time.Now().UTC().UnixNano() / 1e6
	loginResult := make(map[string]interface{})
	if err := c.ShouldBindJSON(&JSONInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"details": "Invalid Input",
		})
		return
	}
	loginResult = db.UserLoginMongo(sessionid, JSONInput) //返回值应当已经预处理完，可以直接放redis
	db.SessionStartRedis(sessionid, loginResult)
	c.SetCookie("SESSIONID", sessionid, COOKIE_EXPIRE_TIME, "/", COOKIE_URL, true, false)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"details":   "登陆成功",
		"sessionid": sessionid,
	})
}
func RegisterHandler(c *gin.Context) {
	var JSONInput RegisterRequest
	ip, _ := c.Request.Header["X-Real-Ip"]
	if err := c.ShouldBindJSON(&JSONInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"details": "无效的输入",
		})
		return
	}
	if db.UserExistsMongo(JSONInput.Email) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"details": "邮箱已被使用",
		})
		c.Abort()
		return
	}
	insertedID := db.CreateNewUser(GetRegisterDocument(JSONInput, ip[0]))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"uid":     insertedID,
		"details": "注册成功",
	})
}

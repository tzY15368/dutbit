package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"./db/RedisConnector"
	//"context"
	//"fmt"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/wonderivan/logger"

	"github.com/garyburd/redigo/redis"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterDocument struct {
	Username        string                 `json:"username" binding:"required"`
	Email           string                 `json:"email" binding:"required"`
	Password        string                 `json:"password" binding:"required"`
	Role            int                    `json:"role"`
	Site            map[string]interface{} `json:"site"`
	Created_at      int64                  `json:"created_at"`
	Ip              string                 `json:"ip"`
	Last_login_time int64                  `json:"last_login_time"`
	Last_login_ip   string                 `json:"last_login_ip"`
	Confirmation    map[string]interface{} `json:"confirmation"`
}

var redis_conn redis.Conn
var mongo_cli *mongo.Client
var mongo_conn *mongo.Database
var UserCollection *mongo.Collection

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func GetSessionId(addr string) string {
	rand.Seed(time.Now().UnixNano())
	return MD5(addr + string(time.Now().UTC().UnixNano()/1e6) + string(rand.Intn(1000)))
}
func ErrorHandler(err error, ErrorType string) {
	if err != nil {
		logger.Error(ErrorType)
		panic(err)
	}
}
func GetRegisterDocument(RegRequest RegisterRequest, UserIp string) RegisterDocument {
	now := time.Now().UTC().UnixNano() / 1e6
	return RegisterDocument{
		Username:        RegRequest.Username,
		Email:           RegRequest.Email,
		Password:        RegRequest.Password,
		Role:            0,
		Site:            make(map[string]interface{}),
		Created_at:      now,
		Last_login_time: now,
		Ip:              UserIp,
		Last_login_ip:   UserIp,
		Confirmation:    make(map[string]interface{}),
	}
}
func session_start(c *gin.Context, sessioninfo map[string]interface{}) string {
	IpArray, _ := c.Request.Header["X-Real-Ip"]
	Ip := IpArray[0]
	sessionid := GetSessionId(Ip)
	//_, err :=
	redis_conn.Do("hmset", redis.Args{}.Add(sessionid).AddFlat(sessioninfo)...)
	//ErrorHandler(err, "hmset error(redis)")
	c.SetCookie("SESSIONID", sessionid, 2592000, "/", ".dutbit.com", true, false)
	return sessionid
}
func RedisInit() {
	cli, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("")) //add password here
	ErrorHandler(err, "redis connection error")
	logger.Info("redis connected")
	redis_conn = cli
}
func MongoInit() {
	clientOptions := options.Client().ApplyURI("mongodb://admin:@localhost:27017/?authSource=admin") //add password here
	client, err := mongo.Connect(context.TODO(), clientOptions)
	ErrorHandler(err, "mongodb connection error")
	err = client.Ping(context.TODO(), nil)
	ErrorHandler(err, "mongodb ping error")
	logger.Info("mongodb connected")
	mongo_cli = client
	mongo_conn = client.Database("userservice2")
	UserCollection = mongo_conn.Collection("users")
}
func AuthRequired(c *gin.Context) {
	sessionid, err := c.Cookie("SESSIONID")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://www.dutbit.com/userservice/index?authrequired")
		c.Abort()
		logger.Info("no cookie. redirecting")
		return
	}
	res, _ := redis.Bool(redis_conn.Do("exists", sessionid))
	if res {
		logger.Info("session ok")
		c.Next()
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "https://www.dutbit.com/userservice/index?authrequired")
		c.Abort()
		return
	}
}
func logout_handler(c *gin.Context) {
	sessionid, _ := c.Cookie("SESSIONID")
	result, err := redis.Int(redis_conn.Do("del", sessionid))
	ErrorHandler(err, "logout error")
	logger.Info("logging out: deleted", result, err)
	c.SetCookie("SESSIONID", "removed", -1, "/", ".dutbit.com", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "https://www.dutbit.com/wp20/index.php?logout")
	return
}
func userinfo_handler(c *gin.Context) {
	sessionid, _ := c.Cookie("SESSIONID")
	logger.Info("cookie:", sessionid)
	result, err := redis.StringMap(redis_conn.Do("hgetall", sessionid))
	ErrorHandler(err, "userinfo handler (redis)error")
	c.JSON(200, result)
}
func login_handler(c *gin.Context) {
	var JSONInput LoginRequest
	original := make(map[string]interface{})
	if err := c.ShouldBindJSON(&JSONInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"details": "无效的输入",
		})
		return
	}
	Ip, _ := c.Request.Header["X-Real-Ip"]
	now := time.Now().UTC().UnixNano() / 1e6
	logger.Info("login: ", Ip)
	filter := bson.M{"email": JSONInput.Email, "password": JSONInput.Password}
	updateField := make(map[string]interface{})
	updateField["last_login_ip"] = Ip[0]
	updateField["last_login_time"] = now
	update := bson.D{
		{"$set", updateField},
	}
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&original)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"details": "邮箱或密码错误",
		})
		c.Abort()
		return
	}
	siteJSON, _ := json.Marshal(original["site"])
	original["site"] = siteJSON
	//objId,err:= primitive.ObjectIDFromHex(original["_id"])
	//ErrorHandler(err,"fatal")

	logger.Info(original["_id"].(primitive.ObjectID).Hex())
	//c.JSON(200, original)
	//return
	original["_id"] = original["_id"].(primitive.ObjectID).Hex()
	UpdateResult, err := UserCollection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	ErrorHandler(err, "userinfo update error(mongodb)")
	if UpdateResult.ModifiedCount != 1 || UpdateResult.MatchedCount != 1 {
		logger.Info(1)
	}
	/*
		sessionid := GetSessionId(Ip[0])
		_, erro := redis_conn.Do("hmset", redis.Args{}.Add(sessionid).AddFlat(original)...)

		c.SetCookie("SESSIONID", sessionid, 2592000, "/", ".dutbit.com", true, false)
	*/
	sessionid := session_start(c, original)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"details":   "登陆成功",
		"sessionid": sessionid,
	})
}
func register_handler(c *gin.Context) {
	var JSONInput RegisterRequest
	original := make(map[string]interface{})
	if err := c.ShouldBindJSON(&JSONInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"details": "无效的输入",
		})
		return
	}
	filter := bson.M{"email": JSONInput.Email}
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&original)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"details": "邮箱已被使用",
		})
		c.Abort()
		return
	}
	Ip, _ := c.Request.Header["X-Real-Ip"]
	logger.Info("register: ", Ip)
	insertResult, err := UserCollection.InsertOne(context.TODO(), GetRegisterDocument(JSONInput, Ip[0]))
	ErrorHandler(err, "registration failure(mongodb)")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"uid":     insertResult.InsertedID,
		"details": "注册成功",
	})
}
func main() {
	RedisInit()
	MongoInit()
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	g := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	g.LoadHTMLGlob(path.Join(cwd, "templates/*"))
	g.GET("/userservice/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "https://www.dutbit.com/userservice/index?authrequired")
	})
	g.GET("/userservice/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	g.GET("/userservice/home", AuthRequired, func(c *gin.Context) { c.HTML(http.StatusOK, "home.html", nil) })
	g.GET("/userservice/logout", AuthRequired, logout_handler)
	g.GET("/userservice/sessionexists", func(c *gin.Context) {
		v, _ := c.Cookie("SESSIONID")
		result, _ := redis.Bool(redis_conn.Do("exists", v))
		logger.Warn(result == true)
	})
	g.GET("/userservice/v1/userinfo", AuthRequired, userinfo_handler)
	g.POST("/userservice/v1/login", login_handler)
	g.POST("/userservice/v1/register", register_handler)
	g.Run("127.0.0.1:8810")
	defer mongo_cli.Disconnect(context.TODO())
	defer redis_conn.Close()
}

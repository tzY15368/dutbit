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
	"encoding/json"
	//"fmt"
	"math/rand"
	"net/http"

	//"reflect"
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
type UpdateRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
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
func JSONToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}
	return tempMap
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

/*
* generates a sessionid, puts it in redis, and starts a goroutine to sync sessionid in mongodb??
 */
func session_start(c *gin.Context, sessioninfo map[string]interface{}) string {
	IpArray, _ := c.Request.Header["X-Real-Ip"]
	Ip := IpArray[0]
	sessionid := GetSessionId(Ip)
	//_, err :=
	redis_conn.Do("hmset", redis.Args{}.Add(sessionid).AddFlat(sessioninfo)...)
	redis_conn.Do("EXPIRE", sessionid, 2592000)
	//ErrorHandler(err, "hmset error(redis)")
	c.SetCookie("SESSIONID", sessionid, 2592000, "/", ".dutbit.com", true, false)
	return sessionid
}
func RedisInit() {
	cli, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("Bit_redis_123"))
	ErrorHandler(err, "redis connection error")
	logger.Info("redis connected")
	redis_conn = cli
}
func MongoInit() {
	clientOptions := options.Client().ApplyURI("mongodb://admin:Bit_root_123@localhost:27017/?authSource=admin")
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
func send(sessionid []string, sessioninfo map[string]interface{}) {
	logger.Info("len of send:", len(sessionid))
	logger.Info("sessioninfo:", sessioninfo)
	for i := 0; i < len(sessionid); i++ {
		if _, err := redis_conn.Do("hmset", redis.Args{}.Add(sessionid[i]).AddFlat(sessioninfo)...); err != nil {
			ErrorHandler(err, "err")
		}
		if _, err := redis_conn.Do("EXPIRE", sessionid[i], 2592000); err != nil {
			ErrorHandler(err, "err")
		}
		logger.Error("redis send:", i)
	}
	if err := redis_conn.Flush(); err != nil {
		ErrorHandler(err, "err")
	}
}

func receive(sessionid []string) {
	logger.Info("len of receive:", len(sessionid))
	for i := 0; i < len(sessionid); i++ {
		_, err := redis_conn.Receive()
		logger.Warn("redis receive", i)
		ErrorHandler(err, "err")
	}
}
func sync_to_redis(sessionid []string, sessioninfo map[string]interface{}) {
	send(sessionid, sessioninfo) //go??
	//receive(sessionid) it's redis.do, not redis.send, no neeed to recv()
}
func userinfo_get_handler(c *gin.Context) {
	sessionid, _ := c.Cookie("SESSIONID")
	logger.Info("cookie:", sessionid)
	result, err := redis.StringMap(redis_conn.Do("hgetall", sessionid))
	ErrorHandler(err, "userinfo handler (redis)error")
	result["password"] = ""
	result["session"] = ""
	c.JSON(200, result)
}
func userinfo_put_handler(c *gin.Context) {
	var JSONInput UpdateRequest
	if err := c.ShouldBindJSON(&JSONInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"details": "Invalid Input" + err.Error(),
		})
		return
	}
	//logger.Info(JSONInput)
	sessionid, _ := c.Cookie("SESSIONID")
	uid, err := redis.String(redis_conn.Do("hget", sessionid, "_id"))
	ErrorHandler(err, "userinfo handler (redis)error")
	ObjId, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.M{"_id": ObjId}
	logger.Info("objid:", ObjId, uid)
	query_map := bson.M{}

	query_map["username"] = JSONInput.Username
	query_map["email"] = JSONInput.Email
	query_map["updated_at"] = time.Now().UTC().UnixNano() / 1e6
	if JSONInput.OldPassword != "" && JSONInput.NewPassword != "" {
		logger.Warn("changing password")
		query_map["password"] = JSONInput.NewPassword
		filter["password"] = JSONInput.OldPassword
	} else {
		logger.Warn("NOT changing passsword")
	}

	update := bson.D{
		{"$set", query_map},
	}
	original := make(map[string]interface{})
	opts := options.FindOneAndUpdate().SetReturnDocument(options.ReturnDocument(1)) //我太强了 蒙对了
	UserCollection.FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		opts,
	).Decode(&original)
	if len(original) == 0 {
		c.JSON(200, gin.H{
			"success": false,
			"details": "Not Found.",
		})
		return
	}
	/*
		updates sessioninfo to all existing redis sessionids
	*/
	sessions := original["session"]
	siteJSON, _ := json.Marshal(original["site"])
	original["site"] = siteJSON
	original["_id"] = original["_id"].(primitive.ObjectID).Hex()
	original["session"] = ""
	var sessionids []string
	for _, v := range sessions.(primitive.A) {
		iv := v.(map[string]interface{})
		sessionids = append(sessionids, iv["sessionid"].(string))
	}
	sync_to_redis(sessionids, original)
	logger.Warn(original)
	c.JSON(200, gin.H{
		"success": true,
		"details": "Personal Info Update Successful",
	})
	return
}
func login_handler(c *gin.Context) {
	var JSONInput LoginRequest
	original := make(map[string]interface{})
	if err := c.ShouldBindJSON(&JSONInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"details": "Invalid Input",
		})
		return
	}
	Ip, _ := c.Request.Header["X-Real-Ip"]
	now := time.Now().UTC().UnixNano() / 1e6
	logger.Info("login: ", Ip)
	filter := bson.M{"email": JSONInput.Email, "password": JSONInput.Password}
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
	original["_id"] = original["_id"].(primitive.ObjectID).Hex()
	sessionid := session_start(c, original)
	sessionField := make(map[string]interface{})
	sessionField["sessionid"] = sessionid
	sessionField["expireAt"] = time.Now().UTC().UnixNano()/1e6 + 2592000
	updateField := make(map[string]interface{})
	updateField["last_login_ip"] = Ip[0]
	updateField["last_login_time"] = now
	update := bson.D{
		{"$set", updateField},
		{"$addToSet", bson.M{
			"session": sessionField,
		},
		},
	} // change this to findone and update
	UpdateResult, err := UserCollection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	ErrorHandler(err, "userinfo update error(mongodb)")
	if UpdateResult.ModifiedCount != 1 || UpdateResult.MatchedCount != 1 {
		logger.Info(1)
	}

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
	g.GET("/userservice/v1/userinfo", AuthRequired, userinfo_get_handler)
	g.PUT("/userservice/v1/userinfo", AuthRequired, userinfo_put_handler)
	g.POST("/userservice/v1/login", login_handler)
	g.POST("/userservice/v1/register", register_handler)
	g.Run("127.0.0.1:8810")
	defer mongo_cli.Disconnect(context.TODO())
	defer redis_conn.Close()
}

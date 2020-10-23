package db

import (
	//"io/ioutil"
	//"encoding/json"
	"context"
	. "models"
	. "utils"

	//"github.com/garyburd/redigo/redis"
	"github.com/wonderivan/logger"
	"go.mongodb.org/mongo-driver/bson"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongo_cli *mongo.Client
var mongo_conn *mongo.Database
var UserCollection *mongo.Collection

func init() {
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

func DeleteSessionMongo(sessionid string) {

}
func UserExistsMongo(email string) bool {
	filter := bson.M{"email": email}
	err := UserCollection.FindOne(context.TODO(), filter)
	if err == nil {
		return true
	}
	return false
}
func CreateNewUser(newUser RegisterDocument) string {
	result, err := UserCollection.InsertOne(context.TODO(), newUser)
	ErrorHandler(err, "create new user failed")
	return result.InsertedID.(string)
}
func UserLoginMongo(sessionid string, userinfo LoginRequest) map[string]string {

}

var Bruh = 2

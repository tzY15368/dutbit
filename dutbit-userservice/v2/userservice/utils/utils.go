package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	. "models"
	"time"

	"github.com/wonderivan/logger"
)

func JSONToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}
	return tempMap
}
func GetNow() int32 {
	return time.Now().UTC().UnixNano() / 1e6
}
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func GetSessionId(addr string) string {
	rand.Seed(time.Now().UnixNano())
	return MD5(addr + string(GetNow()) + string(rand.Intn(1000)))
}
func ErrorHandler(err error, ErrorType string) {
	if err != nil {
		logger.Error(ErrorType)
		panic(err)
	}
}
func GetRegisterDocument(RegRequest RegisterRequest, UserIp string) RegisterDocument {
	now := GetNow()
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

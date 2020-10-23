package main

import (
	"errors"
	"handlers"
	"models"
	. "utils"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

func main() {
	err := errors.New("bruh")
	ErrorHandler(err, "Bbb")
	a := models.Test
	logger.Info(a)
	g := gin.Default()
	g.GET("/userservice/v2/bbb", handlers.LogoutHandler)
	g.GET("/userservice/v2/admin", handlers.LoginRequired, handlers.AdminRequired, handlers.AdminHome)
	g.GET("/userservice/v2/logout", handlers.LoginRequired, handlers.LogoutHandler)
	g.Run("127.0.0.1:8870")
}

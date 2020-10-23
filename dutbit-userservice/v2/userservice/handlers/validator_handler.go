package handlers

import (
	"db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

func LoginRequired(c *gin.Context) {
	//HOME_URL := "https://www.dutbit.com/userservice/home"
	LOGIN_URL := "https://www.dutbit.com/userservice"
	sessionid, err := c.Cookie("SESSIONID")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, LOGIN_URL)
		c.Abort()
		return
	}
	if !db.LoggedIn(sessionid) {
		c.Redirect(http.StatusTemporaryRedirect, LOGIN_URL)
		c.Abort()
		return
	}
	c.Next()
}
func AdminRequired(c *gin.Context) {

	HOME_URL := "https://www.dutbit.com/userservice/home"
	//LOGIN_URL := "https://www.dutbit.com/userservice"
	sessionid, _ := c.Cookie("SESSIONID")
	if !db.IsAdmin(sessionid) {
		c.Redirect(http.StatusTemporaryRedirect, HOME_URL)
		c.Abort()
		return
	}
	logger.Info("Is admin")
	c.Next()
}

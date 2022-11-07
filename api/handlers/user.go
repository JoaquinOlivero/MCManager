package handler

import (
	"MCManager/config"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

func Login(c *gin.Context) {
	type Body struct {
		Password string `json:"password" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	// Get settings
	settings := config.GetValues()
	if settings.Password != body.Password {
		c.String(400, "Wrong password")
		return
	}

	store := ginsession.FromContext(c)
	store.Set("id", store.SessionID())
	store.Save()

	c.Status(200)
}

func Logout(c *gin.Context) {
	store := ginsession.FromContext(c)

	store.Flush()
	store.Save()

	c.Status(200)
}

func ChangePassword(c *gin.Context) {
	type Body struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	// Get settings
	settings := config.GetValues()

	if settings.Password != body.OldPassword {
		c.String(400, "Wrong password")
		return
	}

	// Save new password
	settings.Password = body.NewPassword

	newSettings, err := json.Marshal(settings)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}
	err = ioutil.WriteFile("./config.json", newSettings, 0644)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	c.Status(200)
}

func CheckSession(c *gin.Context) {
	store := ginsession.FromContext(c)
	storeSessionId := store.SessionID()
	cookieSessionId, isExists := store.Get("id")

	if isExists {
		if cookieSessionId != storeSessionId {
			c.Status(401)
			return
		}
	} else {
		c.Status(401)
		return
	}

	c.Status(200)
}

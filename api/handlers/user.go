package handler

import (
	"database/sql"

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

	// Get password from DB
	var dbPassword string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT password FROM settings WHERE id=?", 0)
	err = row.Scan(&dbPassword)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Check if passwords match.
	if dbPassword != body.Password {
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

	// Get password from DB
	var dbPassword string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT password FROM settings WHERE id=?", 0)
	err = row.Scan(&dbPassword)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Check if passwords match.
	if dbPassword != body.OldPassword {
		c.String(400, "Wrong password")
		return
	}

	db.Close()

	// Save new password
	_, err = db.Exec("UPDATE settings SET password=? WHERE id=?", body.NewPassword, 0)
	if err != nil {
		c.String(500, err.Error())
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

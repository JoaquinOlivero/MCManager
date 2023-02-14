package handler

import (
	"database/sql"
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
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	// Get password from DB
	var dbPassword string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT password FROM settings WHERE id=?", 0)
	err = row.Scan(&dbPassword)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Check if passwords match.
	if dbPassword != body.Password {
		log.Println("Wrong password")
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

func SetPassword(c *gin.Context) {
	type Body struct {
		Password string `json:"password" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	// Check if password has already been set.
	var setPassword int
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT setPassword FROM settings WHERE id=?", 0)
	err = row.Scan(&setPassword)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	if setPassword == 1 {
		log.Println("Password has already been set.")
		c.String(401, "Password has already been set.")
		return
	}

	// Update password in database and set "setPassword" to 1 (true).
	_, err = db.Exec("UPDATE settings SET password=?, setPassword=? WHERE id=?", body.Password, 1, 0)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
	}

	db.Close()

	// Sign in user.
	store := ginsession.FromContext(c)
	store.Set("id", store.SessionID())
	store.Save()

	c.Status(200)
}

func CheckSetPassword(c *gin.Context) {
	// Check if password has already been set
	var setPassword int
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT setPassword FROM settings WHERE id=?", 0)
	err = row.Scan(&setPassword)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	if setPassword == 0 {
		c.Status(400)
		return
	}

	db.Close()

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
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	// Get password from DB
	var dbPassword string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT password FROM settings WHERE id=?", 0)
	err = row.Scan(&dbPassword)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	// Check if passwords match.
	if dbPassword != body.OldPassword {
		c.String(400, "Wrong password")
		return
	}

	// Save new password
	_, err = db.Exec("UPDATE settings SET password=? WHERE id=?", body.NewPassword, 0)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

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

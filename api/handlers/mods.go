package handler

import (
	"database/sql"
	"fmt"
	"log"

	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Mods(c *gin.Context) {
	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Set minecraft mods directory path
	modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)
	type mods struct {
		FileName string `json:"fileName"`
		// ModId    string `json:"modId"`
		// Version  string `json:"version"`
	}

	// Walk through "mods" directory and get the files.
	var modsArr []mods

	walkFunc := func(root string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() != "mods" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			modsArr = append(modsArr, mods{FileName: info.Name()})
		}

		return nil
	}

	err = filepath.WalkDir(modsDirectory, walkFunc)
	if err != nil {
		log.Println(err)
		return
	}

	if len(modsArr) == 0 {
		log.Println(err)
		c.Status(204)
		return
	}

	c.JSON(200, modsArr)
}

func UploadMods(c *gin.Context) {
	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Set minecraft mods directory path
	modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	files := form.File["mods"]
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		modPath := fmt.Sprintf("%v/%v", modsDirectory, filename)
		if err := c.SaveUploadedFile(file, modPath); err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"filesUploaded": len(files)})

}

func RemoveMods(c *gin.Context) {
	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Set minecraft mods directory path
	modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)
	// Binding from JSON
	type Body struct {
		ModsList []string `json:"mods" binding:"required"`
	}

	// Bind request body
	var body Body
	err = c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
	}

	// Loop through mods and remove them from the mods directory.
	for _, mod := range body.ModsList {
		modPath := fmt.Sprintf("%v/%v", modsDirectory, mod)

		err := os.Remove(modPath)
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": err})
		}
	}

	c.Status(200)
}

package handler

import (
	"MCManager/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	// Get minecraft directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
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

	// set full filepath
	fullFilePath := fmt.Sprintf("%v%v", minecraftDirectory, c.Query("filepath"))

	// get file extension.
	fileExtension := filepath.Ext(fullFilePath)

	// check if the requested is a backup file with extension bak. If so, then send the actual file format and not the extension. For example alexsmobs.toml.bak --> send ".toml" as the file format.
	if fileExtension == ".bak" {
		fileFormat := filepath.Ext(strings.TrimSuffix(filepath.Base(fullFilePath), ".bak"))
		fileContent, err := utils.FileData(fullFilePath, fileFormat)
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, fileContent)
		return
	} else {
		fileContent, err := utils.FileData(fullFilePath, fileExtension)
		if err != nil {
			log.Println(err)
			c.String(400, err.Error())
			return
		}

		c.JSON(200, fileContent)
		return
	}
}

func SaveFile(c *gin.Context) {
	// Binding the JSON request body
	type Body struct {
		FilePath    string `json:"filepath" binding:"required"`
		FileContent string `json:"file_content" binding:"required"`
	}

	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get minecraft directory from db.
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

	// Set full filepath
	fullFilePath := fmt.Sprintf("%v%v", minecraftDirectory, body.FilePath)

	// Write to file
	err = os.WriteFile(fullFilePath, []byte(body.FileContent), 0660)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.Status(200)
}

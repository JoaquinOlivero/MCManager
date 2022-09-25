package handler

import (
	"MCManager/config"
	"MCManager/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {

	// Get settings
	settings := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := settings.MinecraftDirectory

	// set full filepath
	fullFilePath := fmt.Sprintf("%v%v", minecraftDirectory, c.Query("filepath"))

	// get file extension.
	fileExtension := filepath.Ext(fullFilePath)

	// check if the requested is a backup file with extension bak. If so, then send the actual file format and not the extension. For example alexsmobs.toml.bak --> send ".toml" as the file format.
	if fileExtension == ".bak" {
		fileFormat := filepath.Ext(strings.TrimSuffix(filepath.Base(fullFilePath), ".bak"))
		fileContent, err := utils.FileData(fullFilePath, fileFormat)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, fileContent)
		return
	} else {
		fileContent, err := utils.FileData(fullFilePath, fileExtension)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
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
		// fmt.Println(err) log
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get settings
	settings := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := settings.MinecraftDirectory

	// Set full filepath
	fullFilePath := fmt.Sprintf("%v%v", minecraftDirectory, body.FilePath)

	// Write to file
	err = os.WriteFile(fullFilePath, []byte(body.FileContent), 0660)
	if err != nil {
		// fmt.Println(err) log
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.Status(200)
}

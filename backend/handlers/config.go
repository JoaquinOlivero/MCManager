package handler

import (
	"MCManager/config"
	"MCManager/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ConfigFiles(c *gin.Context) {

	// Get settings
	settings := config.GetValues()

	configDir := fmt.Sprintf("%v/config", settings.MinecraftDirectory)

	directoryFiles, err := utils.DirectoryTree(configDir)
	if err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, directoryFiles)

}

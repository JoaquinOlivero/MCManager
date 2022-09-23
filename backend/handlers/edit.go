package handler

import (
	"MCManager/config"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	// Get filepath from query
	configFilePath := c.Query("filepath")

	// Get settings
	settings := config.GetValues()
	// Set minecraft directory path
	minecraftDirectory := settings.MinecraftDirectory

	fullFilePath := fmt.Sprintf("%v%v", minecraftDirectory, configFilePath)

	// read the whole content of file and pass it to file variable, in case of error pass it to err variable
	file, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	// convert the file binary into a string using string
	fileContent := string(file)

	c.JSON(200, gin.H{
		"file_content": fileContent,
	})

}

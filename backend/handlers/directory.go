package handler

import (
	"MCManager/config"
	"MCManager/utils"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetDirectory(c *gin.Context) {
	// Directory name
	name := c.Param("name")

	// Get settings
	settings := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := settings.MinecraftDirectory

	switch name {
	case "world":
		data, err := worldDir(minecraftDirectory)
		if err != nil {
			c.JSON(400, err)
		}
		c.JSON(200, data)
		return
	case "config":
		data, err := configDir(minecraftDirectory)
		if err != nil {
			c.JSON(400, err)
		}
		c.JSON(200, data)
		return
	}

	c.Status(404)
}

func worldDir(minecraftDirectory string) (interface{}, error) {

	// Open and read server.properties file and retrieve the currrent world name --> level-name=world. The current world name is the name of the directory containing all the world files.
	serverPropertiesPath := fmt.Sprintf("%v/server.properties", minecraftDirectory)

	file, err := os.Open(serverPropertiesPath)
	if err != nil {
		fmt.Println(err) // log
		return nil, err
	}

	// Scan file line by line and retrieve the level-name tag.
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var worldDirName string

	for fileScanner.Scan() {
		line := fileScanner.Text()

		// check if current line is level-name="world"
		targetFound := strings.Contains(line, "level-name")
		if targetFound {
			splitLine := strings.Split(line, "=")
			worldDirName = splitLine[1]
			break
		}
	}
	if err = file.Close(); err != nil {
		fmt.Printf("Could not close the file due to this %s error \n", err) // log
	}

	worldDir := fmt.Sprintf("%v/%v", minecraftDirectory, worldDirName)
	directoryFiles, err := utils.DirectoryTree(worldDir)
	if err != nil {
		fmt.Println(err) // log
		return nil, err
	}

	// c.JSON(200, gin.H{"dir": directoryFiles, "world_name": worldDirName})
	return gin.H{"dir": directoryFiles, "world_name": worldDirName}, nil

}

func configDir(minecraftDirectory string) (interface{}, error) {

	// Set config directory
	configDir := fmt.Sprintf("%v/config", minecraftDirectory)

	directoryFiles, err := utils.DirectoryTree(configDir)
	if err != nil {
		fmt.Println(err) // log
		return nil, err
	}

	return directoryFiles, nil
}

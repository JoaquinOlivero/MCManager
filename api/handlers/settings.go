package handler

import (
	"MCManager/config"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	type DockerContainers struct {
		ContainerId   string `json:"container_id"`
		ContainerName string `json:"container_name"`
	}
	// Get settings
	settings := config.GetValues()

	// Get docker containers if there are any.
	var dockerContainers []DockerContainers
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("Unable to create docker client: %s", err)
	} else {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		if err != nil {
			log.Println("Error when retrieving docker container list: ", err)
		}

		for _, container := range containers {
			dockerContainers = append(dockerContainers, DockerContainers{ContainerId: container.ID[:10], ContainerName: container.Names[0]})
		}
	}

	cli.Close()
	c.JSON(200, gin.H{"settings": gin.H{"minecraft_directory": settings.MinecraftDirectory, "run_method": settings.RunMethod, "docker_container_id": settings.DockerContainerId, "start_command": settings.StartCommand, "backup": settings.Backup}, "docker_containers": dockerContainers})

}

func ConnectDocker(c *gin.Context) {

	// Get settings from config.json
	settings := config.GetValues()

	// binding from JSON
	type Body struct {
		ContainerId string `json:"container_id" binding:"required"`
	}
	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// connect to docker container and obtain additional information
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), body.ContainerId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	fmt.Println(containerInfo)

	// Save settings in config.json
	settings.RunMethod = "docker"
	settings.DockerContainerId = body.ContainerId
	settings.MinecraftServerIp = containerInfo.NetworkSettings.IPAddress
	settings.MinecraftDirectory = containerInfo.Mounts[0].Source

	newSettings, err := json.Marshal(settings)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("./config.json", newSettings, 0644)
	if err != nil {
		fmt.Println(err)
	}

	cli.Close()

	// c.Status(200)
	c.JSON(200, containerInfo)

}

func DisconnectDocker(c *gin.Context) {
	// Get settings from config.json
	settings := config.GetValues()

	// Reset docker settings in config.json
	settings.RunMethod = ""
	settings.DockerContainerId = ""
	settings.MinecraftDirectory = ""
	settings.MinecraftServerIp = ""
	newSettings, err := json.Marshal(settings)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("./config.json", newSettings, 0644)
	if err != nil {
		fmt.Println(err)
	}

	c.Status(200)
}

func SaveCommand(c *gin.Context) {
	// Get settings from config.json
	settings := config.GetValues()

	// binding from JSON
	type Body struct {
		MinecraftDirectory string `json:"minecraft_directory" binding:"required"`
		StartCommand       string `json:"start_command" binding:"startswith=java"`
	}
	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	// Check that minecraft_directory exists and contains the "mods", "config", "logs" directories and the server.properties file.
	// First, check whether minecraft_directory path is absolute.
	isAbs := filepath.IsAbs(body.MinecraftDirectory)
	if isAbs {
		// Second, check whether the root minecraft directory exists.
		_, err := os.Stat(body.MinecraftDirectory)
		if err != nil {
			isNotExists := os.IsNotExist(err)
			if isNotExists {
				c.String(400, "Directory does not exist")
				return
			}
			c.String(400, err.Error())
			return
		}
		// Third, check whether the minecraft directory contains all the required subdirectories.
		dirsToCheck := []string{"mods", "config", "logs", "server.properties"}
		var subDirectoriesErrors []string

		for _, dir := range dirsToCheck {
			_, err = os.Stat(body.MinecraftDirectory + "/" + dir)
			if err != nil {
				isNotExists := os.IsNotExist(err)
				if isNotExists {
					subDirectoriesErrors = append(subDirectoriesErrors, dir)
				}
			}
		}

		// Send response with directories that were not found.
		if len(subDirectoriesErrors) > 0 {
			c.String(400, strings.Join(dirsToCheck, ", ")+" not found in "+`"`+body.MinecraftDirectory+`" directory`)
			return
		}

	} else {
		c.String(400, "Path provided is not absolute")
		return
	}

	// Save new settings.
	settings.RunMethod = "command"
	settings.MinecraftServerIp = "localhost"
	settings.MinecraftDirectory = body.MinecraftDirectory
	settings.StartCommand = body.StartCommand
	settings.DockerContainerId = ""
	newSettings, err := json.Marshal(settings)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
	}
	err = ioutil.WriteFile("./config.json", newSettings, 0644)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
	}

	c.Status(200)
}

func BackupOption(c *gin.Context) {
	// binding from JSON
	type Body struct {
		Option string `json:"option" binding:"required"`
		Value  *bool  `json:"value" binding:"required"`
	}
	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	// Get settings from config.json
	settings := config.GetValues()

	switch body.Option {
	case "world":
		settings.Backup.World = *body.Value
	case "mods":
		settings.Backup.Mods = *body.Value
	case "config":
		settings.Backup.Config = *body.Value
	case "server_properties":
		settings.Backup.ServerProperties = *body.Value
	}

	// Save settings
	newSettings, err := json.Marshal(settings)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}
	err = ioutil.WriteFile("./config.json", newSettings, 0644)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	c.Status(200)
}

package handler

import (
	"MCManager/config"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	// Get docker containers
	var dockerContainers []DockerContainers
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	for _, container := range containers {
		dockerContainers = append(dockerContainers, DockerContainers{ContainerId: container.ID[:10], ContainerName: container.Names[0]})
	}
	cli.Close()
	c.JSON(200, gin.H{"settings": settings, "docker_containers": dockerContainers})

}

func ConnectDocker(c *gin.Context) {

	// Get settings from config.json
	settings := config.GetValues()

	// binding from JSON
	type Body struct {
		ContainerId string `json:"container_id"`
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

	c.Status(200)

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

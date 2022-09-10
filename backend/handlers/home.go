package handler

import (
	"MCManager/config"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	// "github.com/xrjr/mcutils/pkg/ping"
)

func GetHomeInfo(c *gin.Context) {
	type ServerInfo struct {
		// MOTD          string `json:"motd"`
		// Favicon       string `json:"favicon"`
		// OnlinePlayers int    `json:"online_players"`
		DockerStatus string `json:"docker_status"`
		DockerHealth string `json:"docker_health"`
	}
	// Get settings
	settings := config.GetValues()

	// pingclient := ping.NewClient(settings.MinecraftServerIp, 25565)

	// // Connect opens the connection, and can raise an error for example if the server is unreachable
	// err := pingclient.Connect()
	// if err != nil {
	// 	c.Status(500)
	// }

	// // Handshake is the base request of ping, the one that displays number of players, MOTD, etc...
	// // If all went well, hs contains a field Properties which contains a golang-usable JSON Object
	// hs, err := pingclient.Handshake()
	// if err != nil {
	// 	c.Status(500)
	// }

	// // Disconnect closes the connection
	// err = pingclient.Disconnect()
	// if err != nil {
	// 	c.JSON(500, err)
	// }

	// connect to docker container and obtain additional information
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), settings.DockerContainerId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	// GET request response to client.
	var serverInfo ServerInfo
	// serverInfo.MOTD = hs.Properties.Infos().Description
	// serverInfo.Favicon = hs.Properties.Infos().Favicon
	// serverInfo.OnlinePlayers = hs.Properties.Infos().Players.Online
	serverInfo.DockerStatus = containerInfo.State.Status
	serverInfo.DockerHealth = containerInfo.State.Health.Status

	cli.Close()
	c.JSON(200, serverInfo)

}

func ControlServer(c *gin.Context) {

	// Get settings
	settings := config.GetValues()

	action := c.Query("action")
	method := c.Query("method")
	switch method {
	case "docker":
		switch action {
		case "start":
			res, err := startDockerContainer(settings)
			if err != nil {
				c.JSON(500, err)
			}
			c.Status(res)
		case "stop":
			res, err := stopDockerContainer(settings)
			if err != nil {
				c.JSON(500, err)
			}
			c.Status(res)
		}

	}
}

func startDockerContainer(settings config.Config) (int, error) {
	// connect with docker container
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return 500, err
	}
	// start container
	err = cli.ContainerStart(context.Background(), settings.DockerContainerId, types.ContainerStartOptions{})
	if err != nil {
		return 500, err
	}
	cli.Close()
	return 200, nil
}

func stopDockerContainer(settings config.Config) (int, error) {
	// connect with docker container
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return 500, err
	}
	// 	Stop container
	err = cli.ContainerStop(context.Background(), settings.DockerContainerId, nil)
	if err != nil {
		fmt.Println(err)
	}
	cli.Close()
	return 200, nil
}

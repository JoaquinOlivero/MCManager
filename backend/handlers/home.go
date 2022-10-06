package handler

import (
	"MCManager/config"
	"MCManager/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/xrjr/mcutils/pkg/ping"
	"github.com/xrjr/mcutils/pkg/rcon"
)

func GetHomeInfo(c *gin.Context) {

	type ServerInfo struct {
		DockerStatus string     `json:"docker_status"`
		DockerHealth string     `json:"docker_health"`
		RconEnabled  bool       `json:"rcon_enabled"`
		RconPort     string     `json:"rcon_port"`
		RconPassword string     `json:"rcon_password"`
		Ping         ping.Infos `json:"ping_data"`
	}
	// Get settings
	settings := config.GetValues()

	// connect to docker container and get required info about the container.
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), settings.DockerContainerId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	// Set the docker container's status and health into "serverInfo" variable.
	var serverInfo ServerInfo
	serverInfo.DockerStatus = containerInfo.State.Status
	serverInfo.DockerHealth = containerInfo.State.Health.Status

	cli.Close() // Close connection to docker container.

	// If the docker container is running ping the minecraft server to get data back from it.
	if serverInfo.DockerStatus == "running" {
		pingclient := ping.NewClient(settings.MinecraftServerIp, 25565)

		// Connect opens the connection, and can raise an error for example if the server is unreachable
		err = pingclient.Connect()
		if err != nil {
			c.Status(500)
			return
		}

		// Handshake is the base request of ping, the one that displays number of players, MOTD, etc...
		// If all went well, hs contains a field Properties which contains a golang-usable JSON Object
		hs, err := pingclient.Handshake()
		if err != nil {
			c.Status(500)
			return
		}

		// Disconnect closes the connection
		err = pingclient.Disconnect()
		if err != nil {
			c.JSON(500, err)
			return
		}

		// Set the data pinged into "serverInfo" variable.
		serverInfo.Ping = hs.Properties.Infos()

		// Check server.properties lines for "enable-rcon", "rcon.port" and "rcon.password" keys and set their values in "serverInfo" variable.
		rcon, err := utils.ServerPropertiesLineValue("enable-rcon")
		if err != nil {
			c.JSON(500, err)
			return
		}

		rconPort, err := utils.ServerPropertiesLineValue("rcon.port")
		if err != nil {
			c.JSON(500, err)
			return
		}

		rconPassword, err := utils.ServerPropertiesLineValue("rcon.password")
		if err != nil {
			c.JSON(500, err)
			return
		}

		rconBool, err := strconv.ParseBool(rcon)
		if err != nil {
			c.JSON(500, err)
			return
		}
		serverInfo.RconEnabled = rconBool
		serverInfo.RconPort = rconPort
		serverInfo.RconPassword = rconPassword
	}

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

func SendRconCommand(c *gin.Context) {
	type Body struct {
		Command  string `json:"rcon_command" binding:"required"`
		Password string `json:"rcon_password" binding:"required"`
		Port     int    `json:"rcon_port" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
	}

	// Get settings
	settings := config.GetValues()

	rconResponse, err := rcon.Rcon(settings.MinecraftServerIp, body.Port, body.Password, body.Command)
	if err != nil {
		c.String(400, err.Error())
	}

	c.String(200, rconResponse)
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

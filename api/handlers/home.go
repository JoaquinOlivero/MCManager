package handler

import (
	"MCManager/config"
	"MCManager/utils"
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/xrjr/mcutils/pkg/ping"
	"github.com/xrjr/mcutils/pkg/rcon"
)

type ServerInfo struct {
	RunMethod     string     `json:"run_method"`
	DockerStatus  string     `json:"docker_status"`
	DockerHealth  string     `json:"docker_health"`
	StartCommand  string     `json:"start_command"`
	CommandStatus string     `json:"command_status"`
	RconEnabled   bool       `json:"rcon_enabled"`
	RconPort      string     `json:"rcon_port"`
	RconPassword  string     `json:"rcon_password"`
	Ping          ping.Infos `json:"ping_data"`
}

func GetHomeInfo(c *gin.Context) {

	// Get settings
	settings := config.GetValues()

	// Initialization of server info variable
	var serverInfo ServerInfo
	// Set running method
	serverInfo.RunMethod = settings.RunMethod

	switch serverInfo.RunMethod {
	case "docker":
		serverInfo, err := dockerContainerInfo(settings, serverInfo)
		if err != nil {
			c.JSON(500, err)
		}
		c.JSON(200, serverInfo)

	case "command":
		serverInfo, err := commandInfo(settings, serverInfo)
		if err != nil {
			c.JSON(500, err)
		}
		c.JSON(200, serverInfo)
	default:
		c.Status(500)
	}
}

func ControlServer(c *gin.Context) {

	// Get settings
	settings := config.GetValues()

	action := c.Query("action")
	method := settings.RunMethod
	switch method {
	case "docker":
		switch action {
		case "start":
			res, err := startDockerContainer(settings)
			if err != nil {
				c.JSON(res, err)
			}
			c.Status(res)
		case "stop":
			res, err := stopDockerContainer(settings)
			if err != nil {
				c.JSON(res, err)
			}
			c.Status(res)
		}
	case "command":
		switch action {
		case "start":
			res, err := startCommand(settings)
			if err != nil {
				c.JSON(res, err.Error())
			}
			c.Status(res)
		case "stop":
			res, err := stopCommand(settings)
			if err != nil {
				c.JSON(res, err.Error())
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

func Backup(c *gin.Context) {
	// Get settings
	settings := config.GetValues()

	// Get world name
	worldName, err := utils.ServerPropertiesLineValue("level-name")
	if err != nil {
		fmt.Println(err)
	}

	// Get current time to measure total file compression time.
	timeStart := time.Now()

	// zip filename
	filename := "minecraft-server-backup-" + time.Now().Format("02-Jan-2006-15:04:05") + ".zip"
	// Backup file
	backupFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer backupFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(backupFile)
	defer w.Close()

	walkFunc := func(absPath string, info fs.DirEntry, err error) error {
		fmt.Printf("Compressing: %#v\n", absPath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(absPath)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(strings.TrimLeft(absPath, "/"))
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}

	// Backup config directory.
	err = filepath.WalkDir(settings.MinecraftDirectory+"/config", walkFunc)
	if err != nil {
		fmt.Println(err)
	}

	// Backup current world directory.
	err = filepath.WalkDir(settings.MinecraftDirectory+"/"+worldName, walkFunc)
	if err != nil {
		fmt.Println(err)
	}

	// Backup mods directory.
	err = filepath.WalkDir(settings.MinecraftDirectory+"/mods", walkFunc)
	if err != nil {
		fmt.Println(err)
	}

	// Backup server.properties file
	err = filepath.WalkDir(settings.MinecraftDirectory+"/server.properties", walkFunc)
	if err != nil {
		fmt.Println(err)
	}

	// Check error on close for both the archive zip and the actual zip file.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = backupFile.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Time it took to compress all files.
	t := time.Now()
	elapsed := t.Sub(timeStart)
	fmt.Printf("Backup is ready to be sent. Compressing all the files took: %v\n", elapsed)

	// Anonymous function encapsulating c.File() that sends the backup file, so that c.File() doesn't end the execution of the function. Thus, allowing the handler function to continue executing and remove the temporary backup file already downloaded by the user.
	func(c *gin.Context) {
		// Set HTTP headers.
		c.Header("Content-Type", c.GetHeader("Content-Type"))
		c.Header("Content-Disposition", "attachment; fileholder="+filename)
		c.File(filename)
	}(c)

	// Remove temporary backup file created.
	os.Remove(filename)
}

func dockerContainerInfo(settings config.Config, serverInfo ServerInfo) (ServerInfo, error) {
	// connect to docker container and get required info about the container.
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return serverInfo, err
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), settings.DockerContainerId)
	if err != nil {
		return serverInfo, err
	}

	// Set the docker container's status and health into "serverInfo" variable.
	serverInfo.DockerStatus = containerInfo.State.Status
	serverInfo.DockerHealth = containerInfo.State.Health.Status

	cli.Close() // Close connection to docker container.

	// Get server-port from server.properties file.
	serverPropertiesPort, err := utils.ServerPropertiesLineValue("server-port")
	if err != nil {
		return serverInfo, err
	}

	serverPort, err := strconv.Atoi(serverPropertiesPort)
	if err != nil {
		return serverInfo, err
	}

	// If the docker container is running ping the minecraft server to get data back from it.
	if serverInfo.DockerStatus == "running" && serverInfo.DockerHealth == "healthy" {
		pingclient := ping.NewClient(settings.MinecraftServerIp, serverPort)

		// Connect opens the connection, and can raise an error for example if the server is unreachable
		err = pingclient.Connect()
		if err != nil {
			return serverInfo, err
		}

		// Handshake is the base request of ping, the one that displays number of players, MOTD, etc...
		// If all went well, hs contains a field Properties which contains a golang-usable JSON Object
		hs, err := pingclient.Handshake()
		if err != nil {
			return serverInfo, err
		}

		// Disconnect closes the connection
		err = pingclient.Disconnect()
		if err != nil {
			return serverInfo, err
		}

		// Set the data pinged into "serverInfo" variable.
		serverInfo.Ping = hs.Properties.Infos()

		// Check server.properties lines for "enable-rcon", "rcon.port" and "rcon.password" keys and set their values in "serverInfo" variable.
		rcon, err := utils.ServerPropertiesLineValue("enable-rcon")
		if err != nil {
			return serverInfo, err
		}

		rconPort, err := utils.ServerPropertiesLineValue("rcon.port")
		if err != nil {
			return serverInfo, err
		}

		rconPassword, err := utils.ServerPropertiesLineValue("rcon.password")
		if err != nil {
			return serverInfo, err
		}

		rconBool, err := strconv.ParseBool(rcon)
		if err != nil {
			return serverInfo, err
		}
		serverInfo.RconEnabled = rconBool
		serverInfo.RconPort = rconPort
		serverInfo.RconPassword = rconPassword
	}

	return serverInfo, nil
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

func commandInfo(settings config.Config, serverInfo ServerInfo) (ServerInfo, error) {

	// Get server-port from server.properties file.
	serverPropertiesPort, err := utils.ServerPropertiesLineValue("server-port")
	if err != nil {
		return serverInfo, err
	}

	serverPort, err := strconv.Atoi(serverPropertiesPort)
	if err != nil {
		return serverInfo, err
	}

	pingclient := ping.NewClient(settings.MinecraftServerIp, serverPort)

	// Connect opens the connection, and can raise an error for example if the server is unreachable
	err = pingclient.Connect()
	// An error means that the server couldn't be pinged. However, this could mean that the server is either starting or it's offline.
	if err != nil {
		// Check if the server process is running. In this block scope if the server is running, it means that the server is starting.
		process, err := os.FindProcess(settings.Pid)
		if err != nil {
			log.Printf("Failed to find process: %s\n", err)
		} else {
			processStatus := process.Signal(syscall.Signal(0))
			// nil means that the server process is not running and it's therefore offline.
			if processStatus != nil {
				serverInfo.CommandStatus = "offline"
			} else {
				// A value means in this block scope that the server process is running and it's likely that the server is in its starting phase.
				serverInfo.CommandStatus = "starting"
			}
		}

		process.Release()
	} else {
		// If the minecraft server was successfully pinged, it is therefore online. And now data can be obtained from the minecraft server.
		// Handshake is the base request of ping, the one that displays number of players, MOTD, etc...
		// If all went well, hs contains a field Properties which contains a golang-usable JSON Object
		hs, err := pingclient.Handshake()
		if err != nil {
			return serverInfo, err
		}

		// Disconnect closes the connection
		err = pingclient.Disconnect()
		if err != nil {
			return serverInfo, err
		}

		// Set the data pinged into "serverInfo" variable.
		serverInfo.Ping = hs.Properties.Infos()

		// Check server.properties lines for "enable-rcon", "rcon.port" and "rcon.password" keys and set their values in "serverInfo" variable.
		rcon, err := utils.ServerPropertiesLineValue("enable-rcon")
		if err != nil {
			return serverInfo, err
		}

		rconPort, err := utils.ServerPropertiesLineValue("rcon.port")
		if err != nil {
			return serverInfo, err
		}

		rconPassword, err := utils.ServerPropertiesLineValue("rcon.password")
		if err != nil {
			return serverInfo, err
		}

		rconBool, err := strconv.ParseBool(rcon)
		if err != nil {
			return serverInfo, err
		}
		serverInfo.RconEnabled = rconBool
		serverInfo.RconPort = rconPort
		serverInfo.RconPassword = rconPassword

		serverInfo.CommandStatus = "online"
	}

	serverInfo.StartCommand = settings.StartCommand

	return serverInfo, nil
}

func startCommand(settings config.Config) (int, error) {
	log.Println("Starting Minecraft Server")
	// Check if the server process is already running.
	process, err := os.FindProcess(settings.Pid)
	if err != nil {
		log.Printf("Failed to find process: %s\n", err)
	} else {
		err := process.Signal(syscall.Signal(0))
		if err == nil {
			err = errors.New("server is already running")
			process.Release()
			return 400, err
		}
	}

	// Get start command from settings and split it.
	args := strings.Split(settings.StartCommand, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = settings.MinecraftDirectory
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	// Execute the command and start the new process
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return 400, err
	}

	// Save the process id to settings.
	settings.Pid = cmd.Process.Pid // new process id.

	newSettings, err := json.Marshal(settings)
	if err != nil {
		log.Println(err)
		return 500, err
	}
	err = ioutil.WriteFile("./config.json", newSettings, 0644)
	if err != nil {
		log.Println(err)
		return 500, err
	}

	process.Release()

	return 200, nil
}

func stopCommand(settings config.Config) (int, error) {

	// Check if the server process is running.
	process, err := os.FindProcess(settings.Pid)
	if err != nil {
		log.Printf("Failed to find process: %s\n", err)
		return 500, err
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			err = errors.New("server is not running")
			return 400, err
		}
	}

	log.Println("Stopping Minecraft Server")
	process.Kill()
	process.Wait()

	return 200, nil
}

package handler

import (
	"MCManager/utils"
	"context"
	"database/sql"
	"errors"

	// "fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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
	Ping          ping.Infos `json:"ping_data"`
}

func GetHomeInfo(c *gin.Context) {
	var (
		serverInfo   ServerInfo
		method       sql.NullString // Run method. "docker" or "command".
		startCommand sql.NullString // Server start cli command.
		containerId  sql.NullString // Docker container id.
		serverIp     sql.NullString // serverIp. Defaults to localhost. However, it can be changed in settings if needed.
		pid          sql.NullInt32  // Last known process id of minecraft server.
	)

	// Query db to retrieve settings data.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.JSON(500, err)
		return
	}

	defer db.Close()

	// Query through settings table to get unique row with id 0, which contains all the settings.
	row := db.QueryRow("SELECT method, containerId, serverIp, startCommand, serverPid FROM settings WHERE id = ?", 0)
	err = row.Scan(&method, &containerId, &serverIp, &startCommand, &pid)
	if err != nil {
		log.Println(err)
		c.JSON(500, err)
		return
	}

	db.Close()

	if method.Valid {

		switch method.String {
		case "docker":
			serverInfo, err := dockerContainerInfo(containerId.String, serverIp.String, serverInfo)
			if err != nil {
				log.Println(err)
				c.JSON(500, err)
				return
			}

			c.JSON(200, serverInfo)
			return

		case "command":
			serverInfo, err := commandInfo(serverIp.String, startCommand.String, int(pid.Int32), serverInfo)
			if err != nil {
				log.Println(err)
				c.JSON(500, err)
				return
			}

			c.JSON(200, serverInfo)
			return

		default:
			c.Status(500)
			return
		}
	}

	c.JSON(200, serverInfo)
}

func ControlServer(c *gin.Context) {
	// Run method (docker or cli command)
	var method, containerId, startCliCommand, directory sql.NullString
	var pid sql.NullInt32 // Last known process id of minecraft server.

	// Database connection and query.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.JSON(500, err)
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT method, containerId, startCommand, directory, serverPid FROM settings WHERE id = ?", 0)
	err = row.Scan(&method, &containerId, &startCliCommand, &directory, &pid)
	if err != nil {
		log.Println(err)
		c.JSON(500, err)
		db.Close()
		return
	}
	db.Close()

	action := c.Query("action")
	switch method.String {
	case "docker":
		switch action {
		case "start":
			res, err := startDockerContainer(containerId.String)
			if err != nil {
				log.Println(err)
				c.JSON(res, err)
				return
			}

			c.Status(res)
		case "stop":
			res, err := stopDockerContainer(containerId.String)
			if err != nil {
				log.Println(err)
				c.JSON(res, err)
				return
			}

			c.Status(res)
			return
		}
	case "command":
		switch action {
		case "start":
			res, err := startCommand(startCliCommand.String, directory.String, int(pid.Int32))
			if err != nil {
				log.Println(err)
				c.JSON(res, err.Error())
				return
			}

			c.Status(res)
			return
		case "stop":
			res, err := stopCommand(int(pid.Int32))
			if err != nil {
				log.Println(err)
				c.JSON(res, err.Error())
				return
			}

			c.Status(res)
			return
		}
	}
}

func SendRconCommand(c *gin.Context) {
	// Check if rcon is enabled.
	rconEnable, err := utils.ServerPropertiesLineValue("enable-rcon")
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	if rconEnable == "false" {
		log.Println("Rcon is not enabled")
		c.String(400, "Rcon is not enabled.")
		return
	}

	type Body struct {
		Command string `json:"rcon_command" binding:"required"`
	}

	// Bind request body
	var body Body
	err = c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	// Get server ip from database.
	var serverIp string

	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	row := db.QueryRow("SELECT serverIp FROM settings WHERE id=?", 0)
	err = row.Scan(&serverIp)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	// Get rcon port and password from server.properties file.
	rconPort, err := utils.ServerPropertiesLineValue("rcon.port")
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	rconPassword, err := utils.ServerPropertiesLineValue("rcon.password")
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	rconPortInt, err := strconv.Atoi(rconPort)
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	rconResponse, err := rcon.Rcon(serverIp, rconPortInt, rconPassword, body.Command)
	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	c.String(200, rconResponse)
}

func dockerContainerInfo(containerId, serverIp string, serverInfo ServerInfo) (ServerInfo, error) {
	// connect to docker container and get required info about the container.
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return serverInfo, err
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), containerId)
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
		pingclient := ping.NewClient(serverIp, serverPort)

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

		rconBool, err := strconv.ParseBool(rcon)
		if err != nil {
			return serverInfo, err
		}
		serverInfo.RconEnabled = rconBool
	}

	serverInfo.RunMethod = "docker"

	return serverInfo, nil
}

func startDockerContainer(containerId string) (int, error) {
	// connect with docker container
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return 500, err
	}

	log.Println("Starting Minecraft server.")

	// start container
	err = cli.ContainerStart(context.Background(), containerId, types.ContainerStartOptions{})
	if err != nil {
		return 500, err
	}

	cli.Close()
	return 200, nil
}

func stopDockerContainer(containerId string) (int, error) {
	// connect with docker container
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return 500, err
	}

	log.Println("Stopping Minecraft server.")

	// 	Stop container
	err = cli.ContainerStop(context.Background(), containerId, nil)
	if err != nil {
		return 500, err
	}

	cli.Close()
	return 200, nil
}

func commandInfo(serverIp, startCommand string, pid int, serverInfo ServerInfo) (ServerInfo, error) {
	// Get server-port from server.properties file.
	serverPropertiesPort, err := utils.ServerPropertiesLineValue("server-port")
	if err != nil {
		return serverInfo, err
	}

	serverPort, err := strconv.Atoi(serverPropertiesPort)
	if err != nil {
		return serverInfo, err
	}

	pingclient := ping.NewClient(serverIp, serverPort)

	// Connect opens the connection, and can raise an error for example if the server is unreachable
	err = pingclient.Connect()
	// An error means that the server couldn't be pinged. However, this could mean that the server is either starting or it's offline.
	if err != nil {
		// Check if the server process is running. In this block scope if the server is running, it means that the server is starting.
		process, err := os.FindProcess(pid)
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

		rconBool, err := strconv.ParseBool(rcon)
		if err != nil {
			return serverInfo, err
		}
		serverInfo.RconEnabled = rconBool

		serverInfo.CommandStatus = "online"
	}

	serverInfo.StartCommand = startCommand
	serverInfo.RunMethod = "command"

	return serverInfo, nil
}

func startCommand(startCommand, directory string, serverPid int) (int, error) {
	log.Println("Starting Minecraft Server")
	// Check if the server process is already running.
	process, err := os.FindProcess(serverPid)
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
	args := strings.Split(startCommand, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = directory
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	// Execute the command and start the new process
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return 400, err
	}

	// Insert new pid into the database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		return 500, err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE settings SET serverPid = ? WHERE id = ?", cmd.Process.Pid, 0)
	if err != nil {
		return 500, err
	}

	db.Close()

	process.Release()

	return 200, nil
}

func stopCommand(serverPid int) (int, error) {
	// Check if the server process is running.
	process, err := os.FindProcess(serverPid)
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

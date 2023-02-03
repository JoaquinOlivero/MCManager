package handler

import (
	"MCManager/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetSettings(c *gin.Context) {
	type DockerContainers struct {
		ContainerId   string `json:"container_id"`
		ContainerName string `json:"container_name"`
	}

	// Get settings from the database.
	var (
		minecraftDirectory sql.NullString
		method             sql.NullString
		containerId        sql.NullString
		startCommand       sql.NullString
		world              int
		mods               int
		config             int
		serverProperties   int
	)

	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	// Query data from "settings" table.
	row := db.QueryRow("SELECT directory, method, containerId, startCommand FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory, &method, &containerId, &startCommand)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Query data from "backup" table.
	row2 := db.QueryRow("SELECT world, mods, config, serverProperties FROM backup WHERE id=?", 0)
	err = row2.Scan(&world, &mods, &config, &serverProperties)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	db.Close()

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
	c.JSON(200, gin.H{"settings": gin.H{"minecraft_directory": minecraftDirectory.String, "run_method": method.String, "docker_container_id": containerId.String, "start_command": startCommand.String, "backup": gin.H{"world": utils.Itob(world), "mods": utils.Itob(mods), "config": utils.Itob(config), "server_properties": utils.Itob(serverProperties)}}, "docker_containers": dockerContainers})

}

func ConnectDocker(c *gin.Context) {
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

	// Save settings to database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
	}

	defer db.Close()

	_, err = db.Exec("UPDATE settings SET method=?, containerId=?, serverIp=?, directory=?, startCommand=? WHERE id=?", "docker", body.ContainerId, containerInfo.NetworkSettings.IPAddress, containerInfo.Mounts[0].Source, nil, 0)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
	}

	db.Close()

	cli.Close()

	c.JSON(200, containerInfo)

}

func DisconnectDocker(c *gin.Context) {
	// Reset values in database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("UPDATE settings SET method=?, containerId=?, serverIp=?, directory=?, startCommand=? WHERE id=?", nil, nil, "localhost", nil, nil, 0)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
	}

	db.Close()

	c.Status(200)
}

func SaveCommand(c *gin.Context) {
	// binding from JSON
	type Body struct {
		MinecraftDirectory string `json:"minecraft_directory" binding:"required"`
		StartCommand       string `json:"start_command" binding:"startswith=java"`
	}
	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		_, isValidationErr := err.(validator.ValidationErrors)
		if isValidationErr {
			for _, validationErr := range err.(validator.ValidationErrors) {
				tag := validationErr.Tag()
				if tag == "startswith" {
					err := `The command needs to start with: "java"`
					c.String(400, err)
					return
				}
			}
		}

		// handles other err type
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

	// Save new settings in the database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		c.String(500, err.Error())
	}

	defer db.Close()

	_, err = db.Exec("UPDATE settings SET method=?, serverIp=?, directory=?, startCommand=?, containerId=? WHERE id=?", "command", "localhost", body.MinecraftDirectory, body.StartCommand, nil, 0)
	if err != nil {
		c.String(500, err.Error())
	}

	db.Close()

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

	// Convert boolean body.Value to int. 1 if true 0 if false.
	var value int

	if *body.Value {
		value = 1
	} else {
		value = 0
	}

	if body.Option == "world" || body.Option == "mods" || body.Option == "config" || body.Option == "serverProperties" {
		err = saveBackupOption(body.Option, value)
		if err != nil {
			log.Println(err)
			c.String(500, err.Error())
		}

		c.Status(200)
		return
	} else {
		c.Status(400)
		return
	}

}

func CheckSettings(c *gin.Context) {
	// Get settings from the database.
	var (
		minecraftDirectory sql.NullString
		method             sql.NullString
	)

	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	// Query data from "settings" table.
	row := db.QueryRow("SELECT directory, method FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory, &method)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	if !minecraftDirectory.Valid && !method.Valid {
		db.Close()
		c.Status(204)
		return
	}

	db.Close()

	c.Status(200)
}

func saveBackupOption(option string, value int) error {
	// Save option to database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		return err
	}

	defer db.Close()

	statement := fmt.Sprintf("UPDATE backup SET %v=%v WHERE id=%v", option, value, 0)

	_, err = db.Exec(statement)
	if err != nil {
		return err
	}

	db.Close()

	return nil
}

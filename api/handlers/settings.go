package handler

import (
	"MCManager/utils"
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
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
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	// Query data from "settings" table.
	row := db.QueryRow("SELECT directory, method, containerId, startCommand FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory, &method, &containerId, &startCommand)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	// Query data from "backup" table.
	row2 := db.QueryRow("SELECT world, mods, config, serverProperties FROM backup WHERE id=?", 0)
	err = row2.Scan(&world, &mods, &config, &serverProperties)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// connect to docker container and obtain additional information
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err})
		return
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), body.ContainerId)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err})
		return
	}

	// Save settings to database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	var minecraftDirectory string

	_, err = os.Stat("/.dockerenv")
	if os.IsNotExist(err) {
		minecraftDirectory = containerInfo.Mounts[0].Source
	} else {
		log.Println("MCManager is running in a docker container.")
		minecraftDirectory = "/mc"
	}

	_, err = db.Exec("UPDATE settings SET method=?, containerId=?, serverIp=?, directory=?, startCommand=? WHERE id=?", "docker", body.ContainerId, containerInfo.NetworkSettings.IPAddress, minecraftDirectory, nil, 0)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
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
		return
	}

	db.Close()

	c.Status(200)
}

func SaveCommand(c *gin.Context) {
	// binding from JSON
	type Body struct {
		MinecraftDirectory string `json:"minecraft_directory" binding:"required"`
		Script             string `json:"script" binding:"required"`
	}
	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(400, err.Error())
		return
	}

	// Check that minecraft_directory exists and contains the "logs" directories and the server.properties file.
	// First, check whether minecraft_directory path is absolute.
	isAbs := filepath.IsAbs(body.MinecraftDirectory)
	if isAbs {
		// Second, check whether the root minecraft directory exists.
		_, err := os.Stat(body.MinecraftDirectory)
		if err != nil {
			isNotExists := os.IsNotExist(err)
			if isNotExists {
				log.Println(err)
				c.String(400, "Directory does not exist")
				return
			}
			log.Println(err)
			c.String(400, err.Error())
			return
		}
		// Third, check whether the minecraft directory contains the logs subdirectory and the server.properties file.
		dirsToCheck := []string{"logs", "server.properties"}
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
			log.Println(strings.Join(dirsToCheck, ", ") + " not found in " + `"` + body.MinecraftDirectory + `" directory`)
			c.String(400, strings.Join(dirsToCheck, ", ")+" not found in "+`"`+body.MinecraftDirectory+`" directory`)
			return
		}

	} else {
		log.Println("Path provided is not absolute")
		c.String(400, "Path provided is not absolute")
		return
	}

	// Get command from the script file. Read file line by line and retrieve the line that starts with the word "java".
	var command string
	file, err := os.Open(filepath.Clean(body.MinecraftDirectory + "/" + body.Script))
	if err != nil {
		log.Println("Error opening file:", err)
		c.String(500, err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "java") {
			command = line
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error scanning file:", err)
	}

	// Save new settings in the database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec("UPDATE settings SET method=?, serverIp=?, directory=?, startCommand=?, containerId=? WHERE id=?", "command", "localhost", body.MinecraftDirectory, command, nil, 0)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	c.Status(200)
}

func ScriptsInDir(c *gin.Context) {
	// binding from JSON
	type Body struct {
		Dir string `json:"directory" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var files []string

	// First, check whether body.Dir path is absolute.
	isAbs := filepath.IsAbs(body.Dir)
	if isAbs {
		// Second, check whether body.Dir exists.
		_, err := os.Stat(body.Dir)
		if err != nil {
			isNotExists := os.IsNotExist(err)
			if isNotExists {
				log.Println(err)
				c.String(400, "Directory does not exist")
				return
			}
			log.Println(err)
			c.String(400, err.Error())
			return
		}

		// Third, find and return all the .sh files inside body.Dir.
		err = filepath.WalkDir(body.Dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Println(err)
				return err
			}

			if filepath.Ext(path) == ".sh" {
				files = append(files, filepath.Base(path))
			}

			return nil
		})

	} else {
		log.Println("Path provided is not absolute")
		c.String(400, "Path provided is not absolute")
		return
	}

	c.JSON(200, gin.H{"files": files})
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
		log.Println(err)
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
			return
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
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	// Query data from "settings" table.
	row := db.QueryRow("SELECT directory, method FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory, &method)
	if err != nil {
		log.Println(err)
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

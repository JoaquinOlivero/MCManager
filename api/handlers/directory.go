package handler

import (
	"MCManager/utils"
	"bufio"
	"compress/gzip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetDirectory(c *gin.Context) {
	// Directory name
	name := c.Param("name")

	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	switch name {
	case "world":
		data, err := worldDir(minecraftDirectory)
		if err != nil {
			log.Println(err)
			c.String(400, err.Error())
			return
		}

		c.JSON(200, data)
		return
	case "config":
		data, err := configDir(minecraftDirectory)
		if err != nil {
			log.Println(err)
			c.String(400, err.Error())
			return
		}

		c.JSON(200, data)
		return
	case "logs":
		data, err := logsDir(minecraftDirectory)
		if err != nil {
			log.Println(err)
			c.String(400, err.Error())
			return
		}

		c.JSON(200, data)
		return
	}

	c.Status(404)
}

func RemoveFiles(c *gin.Context) {
	// Binding from JSON
	type Body struct {
		FileList  []string `json:"files" binding:"required"`
		Directory string   `json:"directory" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Loop through files and remove them from the directory.
	for _, file := range body.FileList {
		filePath := fmt.Sprintf("%v%v%v", minecraftDirectory, body.Directory, file)
		fileStat, err := os.Lstat(filePath)
		if err != nil {
			log.Println(err)
		}

		if fileStat.IsDir() {
			err := os.RemoveAll(filePath)
			if err != nil {
				log.Println(err)
				c.JSON(500, gin.H{"error": err})
				break
			}
		} else {
			err := os.Remove(filePath)
			if err != nil {
				log.Println(err)
				c.JSON(500, gin.H{"error": err})
				break
			}
		}

	}

	c.Status(200)
}

func ExtractLogFiles(c *gin.Context) {
	// Binding from JSON
	type Body struct {
		FileList []string `json:"files" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(500, err)
		return
	}

	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	db.Close()

	// Loop through files and remove them from the directory.
	for _, file := range body.FileList {
		// Check that the file to be extracted has the .gz extension
		if filepath.Ext(file) != ".gz" {
			log.Println(file + " can't be extracted.")
			continue
		}

		filePath := fmt.Sprintf("%v%v%v", minecraftDirectory, "/logs/", file)

		// Check that file exists.
		_, err := os.Lstat(filePath)
		if err != nil {
			var pErr *os.PathError
			if errors.As(err, &pErr) {
				errMessage := fmt.Sprintf(`"%v" %v`, file, pErr.Err)
				log.Println(errMessage)
				c.String(500, errMessage)
			}
			break
		}

		log.Println("Extracting: " + filePath)

		// Open .gz file
		gzFile, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
			c.String(500, err.Error())
			break
		}

		// Create reader to read the file.
		reader, err := gzip.NewReader(gzFile)
		if err != nil {
			log.Println(err)
			c.String(500, err.Error())
			break
		}
		defer reader.Close()

		// Create new file to write the extracted data
		filename := utils.FileNameWithoutExtSliceNotation(file)
		outfile, err := os.Create(minecraftDirectory + "/logs/" + filename)
		if err != nil {
			log.Println(err)
			c.String(500, err.Error())
			break
		}
		defer outfile.Close()

		// Copy the data from the .gz file to the new file
		_, err = io.Copy(outfile, reader)
		if err != nil {
			log.Println(err)
			c.String(500, err.Error())
			break
		}
	}

	c.Status(200)
}

func worldDir(minecraftDirectory string) (interface{}, error) {
	// Open and read server.properties file and retrieve the currrent world name --> level-name=world. The current world name is the name of the directory containing all the world files.
	serverPropertiesPath := fmt.Sprintf("%v/server.properties", minecraftDirectory)

	file, err := os.Open(serverPropertiesPath)
	if err != nil {
		// Check if server.properties file exists.
		noFile := os.IsNotExist(err)
		if noFile {
			err = errors.New("no such file or directory")
		}
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
		log.Printf("Could not close the file due to this %s error \n", err)
	}

	worldDir := fmt.Sprintf("%v/%v", minecraftDirectory, worldDirName)
	directoryFiles, err := utils.DirectoryTree(worldDir)
	if err != nil {
		return nil, err
	}

	return gin.H{"dir": directoryFiles, "world_name": worldDirName}, nil
}

func configDir(minecraftDirectory string) (interface{}, error) {
	// Set config directory
	configDir := fmt.Sprintf("%v/config", minecraftDirectory)

	directoryFiles, err := utils.DirectoryTree(configDir)
	if err != nil {
		return nil, err
	}

	return directoryFiles, nil
}

func logsDir(minecraftDirectory string) (interface{}, error) {
	// Set logs directory
	logsDir := fmt.Sprintf("%v/logs", minecraftDirectory)

	directoryFiles, err := utils.DirectoryTree(logsDir)
	if err != nil {
		return nil, err
	}

	return directoryFiles, nil
}

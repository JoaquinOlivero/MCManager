package utils

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Node represents a node in a directory tree.
type Node struct {
	Name     string  `json:"name"`
	Children []*Node `json:"children"`
	Parent   *Node   `json:"-"`
	Type     string  `json:"type"`
}

func DirectoryTree(dir string) (result *Node, err error) {

	absRoot, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	parents := make(map[string]*Node)
	walkFunc := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			parents[path] = &Node{
				Name:     filepath.Base(path),
				Children: make([]*Node, 0),
				Type:     "dir",
			}
		} else {
			parents[path] = &Node{
				Name: filepath.Base(path),
				Type: "file",
			}
		}
		return nil
	}

	if err = filepath.WalkDir(absRoot, walkFunc); err != nil {
		return nil, err
	}

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root node.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}

	return

}

func FileData(fullFilePath, fileFormat string) (map[string]interface{}, error) {
	// Supported file extensions
	supportedExtensions := [7]string{".toml", ".json", ".json5", ".properties", ".txt", ".cfg", ".log"}

	var isSupported bool

	for i := range supportedExtensions {
		if fileFormat == supportedExtensions[i] {
			isSupported = true
		}
	}

	if !isSupported {
		err := fmt.Sprintf("%v files are not supported.", fileFormat)
		return nil, errors.New(err)
	}

	file, err := os.ReadFile(fullFilePath)
	if err != nil {
		// Check if file does not exist.
		noFile := os.IsNotExist(err)
		if noFile {
			err = errors.New("no such file or directory")
		}
		return nil, err
	}

	// convert the file binary into a string.
	fileContent := string(file)

	m := map[string]interface{}{
		"file_content": fileContent,
		"file_format":  fileFormat,
	}

	return m, nil

}

func ServerPropertiesLineValue(key string) (string, error) {
	// Get Minecraft server files directory from db.
	var minecraftDirectory string
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		return "", err
	}

	defer db.Close()

	row := db.QueryRow("SELECT directory FROM settings WHERE id=?", 0)
	err = row.Scan(&minecraftDirectory)
	if err != nil {
		return "", err
	}

	db.Close()

	// Open and read server.properties file and retrieve the currrent world name --> level-name=world. The current world name is the name of the directory containing all the world files.
	serverPropertiesPath := fmt.Sprintf("%v/server.properties", minecraftDirectory)

	file, err := os.Open(serverPropertiesPath)
	if err != nil {
		fmt.Println(err) // log
		return "", err
	}

	// Scan file line by line and retrieve the level-name tag.
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var keyValue string

	for fileScanner.Scan() {
		line := fileScanner.Text()

		// check if current line is level-name="world"
		targetFound := strings.Contains(line, key)
		if targetFound {
			splitLine := strings.Split(line, "=")
			keyValue = splitLine[1]
			break
		}
	}
	if err = file.Close(); err != nil {
		fmt.Printf("Could not close the file due to this %s error \n", err) // log
		return "", err
	}

	return keyValue, nil
}

func InitializeDb() error {
	// Check if database file exists. If it doesnt exist, create it.
	dbFile, err := os.OpenFile("config.db", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			log.Println("Database file found.")
			return nil
		} else {
			return err
		}
	}
	log.Println("Creating database file.")
	dbFile.Close()

	// Add tables and default data to the recently created database.
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		return err
	}

	defer db.Close()

	// "settings" table.
	_, err = db.Exec("CREATE TABLE settings (ID int DEFAULT 0, directory VARCHAR(255) DEFAULT NULL, serverIp VARCHAR(255) NOT NULL DEFAULT 'localhost', serverPort INT DEFAULT 25565, method VARCHAR(10) DEFAULT NULL, containerId VARCHAR(255) DEFAULT NULL, startCommand VARCHAR(255) DEFAULT NULL, serverPid INT DEFAULT NULL, password VARCHAR(255) DEFAULT 'admin', setPassword BIT DEFAULT 0)")
	if err != nil {
		return err
	}

	// insert data into settings table.
	_, err = db.Exec("INSERT INTO settings (id) VALUES (0)")
	if err != nil {
		return err
	}

	// "backup" table.
	_, err = db.Exec("CREATE TABLE backup (ID int NOT NULL DEFAULT 0, world bit NOT NULL DEFAULT 1, mods bit NOT NULL DEFAULT 1, config bit NOT NULL DEFAULT 1, serverProperties bit NOT NULL DEFAULT 1)")
	if err != nil {
		return err
	}

	// insert data into backup table.
	_, err = db.Exec("INSERT INTO backup VALUES (0,1,1,1,1)")
	if err != nil {
		return err
	}

	db.Close()

	return nil
}

func Itob(i int) bool {
	if i != 0 {
		return i != 0
	}

	return false
}

func FileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

package utils

import (
	"MCManager/config"
	"bufio"
	"errors"
	"fmt"
	"io/fs"
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
	settings := config.GetValues()
	// Open and read server.properties file and retrieve the currrent world name --> level-name=world. The current world name is the name of the directory containing all the world files.
	serverPropertiesPath := fmt.Sprintf("%v/server.properties", settings.MinecraftDirectory)

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

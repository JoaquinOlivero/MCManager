package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	supportedExtensions := [6]string{".toml", ".json", ".json5", ".properties", ".txt", ".cfg"}

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

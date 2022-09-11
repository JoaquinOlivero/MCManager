package utils

import (
	"io/fs"
	"path/filepath"
)

// type File struct {
// 	Name string
// }

// type Folder struct {
// 	Name    string
// 	Files   []string
// 	Folders map[string]*Folder
// }

// func DirectoryTree(dir string) (*Folder, error) {
// 	var tree *Folder
// 	var nodes = map[string]interface{}{}
// 	walk := func(p string, info fs.DirEntry, err error) error {
// 		if info.IsDir() {
// 			nodes[p] = &Folder{path.Base(p), []string{}, map[string]*Folder{}}
// 		} else {
// 			nodes[p] = &File{path.Base(p)}
// 		}
// 		return nil
// 	}

// 	err := filepath.WalkDir(dir, walk)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for key, value := range nodes {
// 		var parentFolder *Folder
// 		if key == dir {
// 			tree = value.(*Folder)
// 			continue
// 		} else {
// 			parentFolder = nodes[path.Dir(key)].(*Folder)
// 		}

// 		switch v := value.(type) {
// 		case *File:
// 			parentFolder.Files = append(parentFolder.Files, v.Name)
// 		case *Folder:
// 			parentFolder.Folders[v.Name] = v
// 		}
// 	}

// 	return tree, nil
// }

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
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}

	return

}

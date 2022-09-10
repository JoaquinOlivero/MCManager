package handler

import (
	"MCManager/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/essentialkaos/go-jar"
	"github.com/gin-gonic/gin"
)

func Mods(c *gin.Context) {
	// config.json values
	config := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := config.MinecraftDirectory
	// Set minecraft mods directory path
	modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)
	type mods struct {
		FileName string `json:"fileName"`
		ModId    string `json:"modId"`
		Version  string `json:"version"`
	}
	var modsArr []mods

	files, err := ioutil.ReadDir(modsDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileName := f.Name()

		// manifest, err := jar.ReadFile(modsDirectory + "/" + fileName)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// // fmt.Println(manifest)
		// // manifest["Implementation-Timestamp"] --> timestamp from when it was created. Useful to check if the mod has new updates.
		// modId := manifest["Implementation-Title"]
		// version := manifest["Implementation-Version"]

		modsArr = append(modsArr, mods{FileName: fileName})

	}

	c.JSON(200, modsArr)
}

func UploadMods(c *gin.Context) {
	// config.json values
	config := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := config.MinecraftDirectory
	// Set minecraft mods directory path
	modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	files := form.File["mods"]
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		modPath := fmt.Sprintf("%v/%v", modsDirectory, filename)
		if err := c.SaveUploadedFile(file, modPath); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"filesUploaded": len(files)})

}

func RemoveMods(c *gin.Context) {
	// config.json values
	config := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := config.MinecraftDirectory
	// Set minecraft mods directory path
	modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)
	// Binding from JSON
	type Body struct {
		ModsList []string `json:"mods" binding:"required"`
	}

	// Bind request body
	var body Body
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Loop through mods and remove them from the mods directory.
	for _, mod := range body.ModsList {
		modPath := fmt.Sprintf("%v/%v", modsDirectory, mod)

		err := os.Remove(modPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
	}

	c.Status(200)
}

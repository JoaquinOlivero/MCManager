package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/essentialkaos/go-jar"
	"github.com/gin-gonic/gin"
)

func Mods(modsDir string) gin.HandlerFunc {
	type mods struct {
		FileName string `json:"fileName"`
		ModId    string `json:"modId"`
		Version  string `json:"version"`
	}
	fn := func(c *gin.Context) {
		var modsArr []mods

		files, err := ioutil.ReadDir(modsDir)
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
	return gin.HandlerFunc(fn)
}

func UploadMods(modsDir string) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		files := form.File["mods"]
		for _, file := range files {
			filename := filepath.Base(file.Filename)
			modPath := fmt.Sprintf("%v/%v", modsDir, filename)
			if err := c.SaveUploadedFile(file, modPath); err != nil {
				c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"filesUploaded": len(files)})
	}

	return gin.HandlerFunc(fn)

}

func RemoveMods(modsDir string) gin.HandlerFunc {
	// Binding from JSON
	type Body struct {
		ModsList []string `json:"mods" binding:"required"`
	}

	fn := func(c *gin.Context) {

		// Bind request body
		var body Body
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// Loop through mods and remove them from the mods directory.
		for _, mod := range body.ModsList {
			modPath := fmt.Sprintf("%v/%v", modsDir, mod)

			err := os.Remove(modPath)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}
		}

		c.Status(200)
	}
	return gin.HandlerFunc(fn)
}

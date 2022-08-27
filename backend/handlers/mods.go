package handler

import (
	"fmt"
	"io/ioutil"
	"log"

	// "github.com/essentialkaos/go-jar"
	"github.com/gin-gonic/gin"
)

func Mods(dir string) gin.HandlerFunc {
	type mods struct {
		FileName string `json:"fileName"`
		ModId    string `json:"modId"`
		Version  string `json:"version"`
	}
	fn := func(c *gin.Context) {
		var modsArr []mods

		modsDirectory := fmt.Sprintf("%v/mods", dir)

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
	return gin.HandlerFunc(fn)
}

package handler

import (
	"MCManager/utils"
	"archive/zip"
	"database/sql"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Backup(c *gin.Context) {
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

	// Get world name
	worldName, err := utils.ServerPropertiesLineValue("level-name")
	if err != nil {
		log.Println(err)
	}

	// Get current time to measure total file compression time.
	timeStart := time.Now()

	// zip filename
	filename := "minecraft-server-backup-" + time.Now().Format("02-Jan-2006-15:04:05") + ".zip"

	// Create backup directory and backup file.
	backupFile, err := os.Create("backup/" + filename)
	if err != nil {
		// if the error is that the "backup" directory does not exist, create it.
		if os.IsNotExist(err) {
			err = os.Mkdir("backup", 0755)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			log.Println(err)
			return
		}

	}
	defer backupFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(backupFile)
	defer w.Close()

	walkFunc := func(absPath string, info fs.DirEntry, err error) error {
		log.Printf("Compressing: %#v\n", absPath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(absPath)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(strings.TrimLeft(absPath, "/"))
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}

	// Backup config directory.
	err = filepath.WalkDir(minecraftDirectory+"/config", walkFunc)
	if err != nil {
		log.Println(err)
	}

	// Backup current world directory.
	err = filepath.WalkDir(minecraftDirectory+"/"+worldName, walkFunc)
	if err != nil {
		log.Println(err)
	}

	// Backup mods directory.
	err = filepath.WalkDir(minecraftDirectory+"/mods", walkFunc)
	if err != nil {
		log.Println(err)
	}

	// Backup server.properties file
	err = filepath.WalkDir(minecraftDirectory+"/server.properties", walkFunc)
	if err != nil {
		log.Println(err)
	}

	// Check error on close for both the archive zip and the actual zip file.
	err = w.Close()
	if err != nil {
		log.Println(err)
	}

	err = backupFile.Close()
	if err != nil {
		log.Println(err)
	}

	// Time it took to compress all files.
	t := time.Now()
	elapsed := t.Sub(timeStart)
	log.Printf("Backup is ready. Compressing all the files took: %v\n", elapsed)

	c.String(200, filename)
}

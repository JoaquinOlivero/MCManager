package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type BackupOptions struct {
	World            bool `json:"world"`
	Mods             bool `json:"mods"`
	Config           bool `json:"config"`
	ServerProperties bool `json:"server_properties"`
}

type Config struct {
	MinecraftDirectory string        `json:"minecraft_directory"`
	MinecraftServerIp  string        `json:"minecraft_server_ip"`
	RunMethod          string        `json:"run_method"`
	DockerContainerId  string        `json:"docker_container_id"`
	StartCommand       string        `json:"start_command"`
	Pid                int           `json:"server_pid"`
	Backup             BackupOptions `json:"backup"`
	Password           string        `json:"password"`
}

func GetValues() Config {
	var config Config

	// Read config.json
	content, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error when opening config file: ", err)

		noConfigFile := os.IsNotExist(err)
		if noConfigFile {
			fmt.Println("Creating config.json file and setting default values")

			// Set default config values.
			config.Backup.World = true
			config.Backup.Mods = true
			config.Backup.Config = true
			config.Backup.ServerProperties = true
			config.Pid = 0
			config.Password = "admin"

			// Marshall config struct into json []byte.
			jsonConfig, _ := json.Marshal(config)
			err := os.WriteFile("./config.json", jsonConfig, 0666)
			if err != nil {
				fmt.Println("Error when creating config file: ", err)
			}
			return config
		}

	}

	// Unmarshall config content
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Error during config.json Unmarshal(): ", err)
	}

	return config

}

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MinecraftDirectory string `json:"minecraft_directory"`
	MinecraftServerIp  string `json:"minecraft_server_ip"`
	RunMethod          string `json:"run_method"`
	DockerContainerId  string `json:"docker_container_id"`
	StartScript        string `json:"start_script"`
	StopScript         string `json:"stop_script"`
}

func GetValues() Config {
	var config Config

	// Read config.json
	content, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error when opening config file: ", err)

		noConfigFile := os.IsNotExist(err)
		if noConfigFile {
			fmt.Println("Creating config.json file")
			// Marshall config struct into json []byte
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

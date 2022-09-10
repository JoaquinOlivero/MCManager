package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	// Read config.json
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error when opening config file: ", err)
	}
	// Unmarshall config content
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Error during config.json Unmarshal(): ", err)
	}

	return config

}

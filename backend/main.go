package main

import (
	"MCManager/config"
	handler "MCManager/handlers"

	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	// "github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

// type Config struct {
// 	MinecraftDirectory string `json:"minecraft_directory"`
// 	RunMethod          string `json:"run_method"`
// 	DockerContainerId  string `json:"docker_container_id"`
// 	StartScript        string `json:"start_script"`
// 	StopScript         string `json:"stop_script"`
// }

func main() {
	os.Setenv("MCMANAGER_HTTP_PROXY_PORT", "5000")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.SetTrustedProxies(nil)
	// r.Use(static.Serve("/", static.LocalFile("../frontend/out/", false)))

	// config.json values
	config := config.GetValues()

	// Set minecraft directory path
	minecraftDirectory := config.MinecraftDirectory

	// Setup route group for the API
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			// write to config.json file
			// config.RunMethod = "docker"

			// content, err := json.Marshal(config)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// err = ioutil.WriteFile("config.json", content, 0644)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// 74694eebc6f3 // start container
			// err = cli.ContainerStart(context.Background(), "74694eebc6f3", types.ContainerStartOptions{})
			// if err != nil {
			// 	fmt.Println(err)
			// }

			// Stop container
			// err = cli.ContainerStop(context.Background(), "74694eebc6f3", nil)
			// if err != nil {
			// 	fmt.Println(err)
			// }
		})
		mods := api.Group("/mods")
		{
			// Set minecraft mods directory path
			modsDirectory := fmt.Sprintf("%v/mods", minecraftDirectory)
			mods.GET("/", handler.Mods(modsDirectory))
			mods.POST("/upload", handler.UploadMods((modsDirectory)))
			mods.POST("/remove", handler.RemoveMods(modsDirectory))
		}

		settings := api.Group("/settings")
		{
			settings.GET("/", handler.GetSettings)
			settings.POST("/connect-docker", handler.ConnectDocker)
			settings.POST("/disconnect-docker", handler.DisconnectDocker)
		}
	}

	port := os.Getenv("MCMANAGER_HTTP_PROXY_PORT")
	fmt.Printf("Server started on port: %v\n", port)
	r.NoRoute(ReverseProxy) //
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("../frontend/out/index.html")
	// })
	r.Run(":" + port)
}

func ReverseProxy(c *gin.Context) {
	remote, _ := url.Parse("http://localhost:3001")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL = c.Request.URL
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

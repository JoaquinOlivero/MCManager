package main

import (
	"MCManager/config"
	handler "MCManager/handlers"

	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"os"

	"github.com/gin-contrib/static"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
)

func main() {
	config.GetValues()

	os.Setenv("MCMANAGER_HTTP_PROXY_PORT", "5000")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(helmet.Default())
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.SetTrustedProxies(nil)
	r.Use(static.Serve("/", static.LocalFile("out/", false)))

	// Setup route group for the API
	api := r.Group("/api")
	{
		api.GET("/", handler.GetHomeInfo)
		api.POST("/", handler.ControlServer)
		api.POST("/rcon", handler.SendRconCommand)
		api.GET("/backup", handler.Backup)

		dir := api.Group("/dir")
		{
			dir.GET("/:name", handler.GetDirectory)
			dir.POST("/remove", handler.RemoveFiles)
		}
		mods := api.Group("/mods")
		{
			mods.GET("/", handler.Mods)
			mods.POST("/upload", handler.UploadMods)
			mods.POST("/remove", handler.RemoveMods)
		}

		settings := api.Group("/settings")
		{
			settings.GET("/", handler.GetSettings)
			settings.POST("/connect-docker", handler.ConnectDocker)
			settings.POST("/disconnect-docker", handler.DisconnectDocker)
		}

		edit := api.Group("/edit")
		{
			edit.GET("/", handler.GetFile)
			edit.POST("/save", handler.SaveFile)
		}
	}

	port := os.Getenv("MCMANAGER_HTTP_PROXY_PORT")
	// port := "5001"
	fmt.Printf("Server started on port: %v\n", port)
	// r.NoRoute(ReverseProxy)
	r.NoRoute(func(c *gin.Context) {
		c.File("out/index.html")
	})
	r.Run(":" + port)
}

func ReverseProxy(c *gin.Context) {
	remote, _ := url.Parse("http://localhost:3002")
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

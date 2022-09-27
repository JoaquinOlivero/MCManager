package main

import (
	handler "MCManager/handlers"

	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	// "github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

func main() {
	os.Setenv("MCMANAGER_HTTP_PROXY_PORT", "5000")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.SetTrustedProxies(nil)
	// r.Use(static.Serve("/", static.LocalFile("../frontend/out/", false)))

	// Setup route group for the API
	api := r.Group("/api")
	{
		api.GET("/", handler.GetHomeInfo)
		api.POST("/", handler.ControlServer)

		dir := api.Group("/dir")
		{
			dir.GET("/:name", handler.GetDirectory)
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

package main

import (
	"MCManager/config"
	handler "MCManager/handlers"
	"MCManager/middleware"

	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	// "os"

	// "github.com/gin-contrib/static"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/go-session/session"
)

func main() {
	config.GetValues()

	// os.Setenv("MCMANAGER_HTTP_PROXY_PORT", "5000")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(helmet.Default())

	r.Use(ginsession.New(session.SetCookieName("MCManager")))
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.SetTrustedProxies(nil)
	// r.Use(static.Serve("/", static.LocalFile("out/", false)))

	// Setup route group for the API
	api := r.Group("/api")
	{
		api.GET("/", middleware.Session, handler.GetHomeInfo)
		api.POST("/", middleware.Session, handler.ControlServer)
		api.POST("/rcon", middleware.Session, handler.SendRconCommand)
		api.GET("/backup", middleware.Session, handler.Backup)
		api.POST("/login", handler.Login)
		api.GET("/check", handler.CheckSession)
		api.GET("/logout", middleware.Session, handler.Logout)

		dir := api.Group("/dir")
		dir.Use(middleware.Session)
		{
			dir.GET("/:name", handler.GetDirectory)
			dir.POST("/remove", handler.RemoveFiles)
		}
		mods := api.Group("/mods")
		mods.Use(middleware.Session)
		{
			mods.GET("/", handler.Mods)
			mods.POST("/upload", handler.UploadMods)
			mods.POST("/remove", handler.RemoveMods)
		}

		settings := api.Group("/settings")
		settings.Use(middleware.Session)
		{
			settings.GET("/", handler.GetSettings)
			settings.POST("/connect-docker", handler.ConnectDocker)
			settings.POST("/disconnect-docker", handler.DisconnectDocker)
			settings.POST("/save-command", handler.SaveCommand)
		}

		edit := api.Group("/edit")
		edit.Use(middleware.Session)
		{
			edit.GET("/", handler.GetFile)
			edit.POST("/save", handler.SaveFile)
		}
	}

	// port := os.Getenv("MCMANAGER_HTTP_PROXY_PORT")
	port := "5001"
	fmt.Printf("Server started on port: %v\n", port)
	r.NoRoute(ReverseProxy)
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("out/index.html")
	// })
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

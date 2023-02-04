package main

import (
	handler "MCManager/handlers"
	"MCManager/middleware"
	"MCManager/utils"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"flag"

	"github.com/gin-contrib/static"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/go-session/session"
)

func main() {
	portFlag := flag.String("p", "", "-p flag defines the port to be used by MCManager. Defaults to 5555.")
	devFlag := flag.Bool("dev", false, "Used to proxy requests to the front-end running in dev mode in port 3002. Therefore, not using the front-end static files in the build folder.")

	flag.Parse()

	port := *portFlag
	dev := *devFlag

	if port == "" {
		port = "5555"
	}

	os.Setenv("MCMANAGER_HTTP_PROXY_PORT", port)

	err := utils.InitializeDb()
	if err != nil {
		log.Println(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(helmet.Default())

	r.Use(ginsession.New(session.SetCookieName("MCManager")))
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.SetTrustedProxies(nil)
	r.Use(static.Serve("/", static.LocalFile("out/", false)))

	// Setup route group for the API
	api := r.Group("/api")
	{
		api.GET("/", middleware.Session, handler.GetHomeInfo)
		api.POST("/", middleware.Session, handler.ControlServer)
		api.POST("/rcon", middleware.Session, handler.SendRconCommand)

		api.POST("/login", handler.Login)
		api.GET("/check", handler.CheckSession)
		api.GET("/logout", middleware.Session, handler.Logout)

		password := api.Group("/password")
		{
			password.POST("/set", handler.SetPassword)
			password.GET("/check", handler.CheckSetPassword)
			password.POST("/change", middleware.Session, handler.ChangePassword)
		}

		dir := api.Group("/dir")
		dir.Use(middleware.Session)
		{
			dir.GET("/:name", handler.GetDirectory)
			dir.POST("/remove", handler.RemoveFiles)
			dir.POST("/extract", handler.ExtractLogFiles)
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
			settings.POST("/save-command", handler.SaveCommand)
			settings.POST("/backup", handler.BackupOption)
			settings.GET("/check", handler.CheckSettings)

			docker := settings.Group("/docker")
			{
				docker.POST("/connect", handler.ConnectDocker)
				docker.POST("/disconnect", handler.DisconnectDocker)
			}

			command := settings.Group("/command")
			{
				command.POST("/save", handler.SaveCommand)
			}
		}

		edit := api.Group("/edit")
		edit.Use(middleware.Session)
		{
			edit.GET("/", handler.GetFile)
			edit.POST("/save", handler.SaveFile)
		}

		backup := api.Group("/backup")
		backup.Use(middleware.Session)
		{
			backup.GET("", middleware.Session, handler.Backup)
			backup.StaticFS("/download", gin.Dir("backup", true))
		}
	}

	// If the dev flag is not used, the back-end will serve the front-end static files in the "out" directory. Otherwise, it will proxy the requests to the port 3002 which is the port that the front-end uses when in dev mode.
	if !dev {
		r.NoRoute(func(c *gin.Context) {
			c.File("out/index.html")
		})

		log.Printf("Server started on port: %v\n", port)
	} else {
		// Proxy to the front-end running in dev mode.
		r.NoRoute(func(c *gin.Context) {
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
		})

		log.Printf("[DEV]Server started on port: %v\n", port)
	}

	r.Run(":" + port)
}

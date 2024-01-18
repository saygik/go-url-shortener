package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	requestid "github.com/thanhhh/gin-requestid"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/saygik/go-url-shortener/internal/config"
	"github.com/saygik/go-url-shortener/internal/http-server/handlers/redirect"
	"github.com/saygik/go-url-shortener/internal/http-server/handlers/url/save"

	// "github.com/saygik/go-url-shortener/internal/http-server/handlers/redirect"
	// "github.com/saygik/go-url-shortener/internal/http-server/handlers/url/save"
	"github.com/saygik/go-url-shortener/internal/lib/logger/sl"
	"github.com/saygik/go-url-shortener/internal/storage/mssql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func main() {
	cfg, err := config.MustLoad(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	log := setupLogger(cfg.Env)
	src, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err", err)
	}
	log.Out = src

	//log.Info("initializing server", " port", cfg.HTTPServer.Port) // Помимо сообщения выведем параметр с адресом

	//	storage, err := sqlite.New(cfg.StoragePath)
	storage, err := mssql.New(mssql.ConnectionParameters{Server: cfg.DBServer, User: cfg.DBUser, Database: cfg.DBName, Password: cfg.DBPassword})

	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	if cfg.Env == envProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(ginlogrus.Logger(log), gin.Recovery())
	r.Use(CORSMiddleware())
	r.Use(requestid.RequestID())

	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFile("/icon.png", "./resources/icon.png")
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	api := r.Group("/api")
	api.POST("/", save.New(log, storage))
	r.GET("/:alias", redirect.New(log, storage))

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	if cfg.Env == envProd {
		gin.SetMode(gin.ReleaseMode)
	}

	http.ListenAndServe(":"+cfg.Port, r)

}

func setupLogger(env string) *logrus.Logger {

	var log = logrus.New()

	// src, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// log.Out = src

	log.Info("------------------Starting programm-------------")

	switch env {
	case envLocal:
		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05"
		customFormatter.FullTimestamp = true
		log.Formatter = customFormatter
	case envProd:
		customFormatter2 := new(logrus.JSONFormatter)
		customFormatter2.TimestampFormat = "2006-01-02 15:04:05"
		log.Formatter = customFormatter2
		log.Debug("logger debug mode enabled")
	}

	return log
}

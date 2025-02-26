package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-url-shortener/internal/lib/constants"
	"github.com/sirupsen/logrus"
	requestid "github.com/thanhhh/gin-requestid"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Server struct {
	Rtr *gin.Engine
}

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

func NewServer(env string, log *logrus.Logger) *Server {
	server := &Server{}
	if env == constants.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}
	server.Rtr = gin.New()
	server.Rtr.Use(ginlogrus.Logger(log), gin.Recovery())
	server.registerMiddlewares()

	return server
}

func (s *Server) registerMiddlewares() {
	s.Rtr.Use(CORSMiddleware())
	s.Rtr.Use(requestid.RequestID())
}

func (s *Server) Start(port int) error {
	err := http.ListenAndServe(":"+fmt.Sprint(port), s.Rtr)
	if err != nil {
		return err
	}
	return nil
}

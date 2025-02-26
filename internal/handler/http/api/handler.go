package api

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	rtr *gin.Engine
	uc  UseCase
	log *logrus.Logger
}

type UseCase interface {
	SaveURL(string, string) (interface{}, error)
	GetURL(string) (string, error)
}

func NewHandler(router *gin.Engine, uc UseCase, log *logrus.Logger) {
	h := &Handler{
		rtr: router,
		uc:  uc,
		log: log,
	}
	h.rtr.StaticFile("/favicon.ico", "./resources/favicon.ico")
	h.rtr.StaticFile("/icon.png", "./resources/icon.png")
	h.rtr.Use(static.Serve("/", static.LocalFile("./public/html", true)))
	h.rtr.NoRoute(h.NoRoute)
	g := h.rtr.Group("/api")
	g.POST("/", h.SaveURL)
	h.rtr.GET("/:alias", h.Redirect)
}

func (h *Handler) NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"Code": "404", "Message": "Not Found"})
	c.Abort()

}

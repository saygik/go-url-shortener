package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	resp "github.com/saygik/go-url-shortener/internal/lib/api/response"
	"github.com/saygik/go-url-shortener/internal/storage"
)

func (h *Handler) Redirect(c *gin.Context) {

	alias := c.Param("alias")
	if alias == "" {
		h.log.Info("Псевдоним пустой")

		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("Неверный запрос"))
		return
	}
	resURL, err := h.uc.GetURL(alias)

	if errors.Is(err, storage.ErrURLNotFound) {
		h.log.Info("url not found", "alias", alias)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("Страница с таким псевдонимом не найдена"))
		return
	}
	if err != nil {
		h.log.Error("failed to get url")

		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("Внутренняя ошибка"))
		return
	}
	h.log.Info("got url ", resURL)
	// redirect to found url
	c.Redirect(http.StatusFound, resURL)
}

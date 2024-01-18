package redirect

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	requestid "github.com/thanhhh/gin-requestid"

	resp "github.com/saygik/go-url-shortener/internal/lib/api/response"
	"github.com/saygik/go-url-shortener/internal/lib/logger/sl"
	"github.com/saygik/go-url-shortener/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *logrus.Logger, urlGetter URLGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.redirect.New"

		log.WithFields(
			logrus.Fields{
				"op":         op,
				"request_id": requestid.GetReqID(c),
			}).Info("")

		alias := c.Param("alias")
		if alias == "" {
			log.Info("Псевдоним пустой")

			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("Неверный запрос"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("Страница с таким псевдонимом не найдена"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("Внутренняя ошибка"))
			return
		}

		log.Info("got url ", resURL)

		// redirect to found url
		c.Redirect(http.StatusFound, resURL)
	}
}

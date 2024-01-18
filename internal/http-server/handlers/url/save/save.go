package save

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/sirupsen/logrus"
	requestid "github.com/thanhhh/gin-requestid"

	resp "github.com/saygik/go-url-shortener/internal/lib/api/response"
	"github.com/saygik/go-url-shortener/internal/lib/random"
	"github.com/saygik/go-url-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config if needed
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *logrus.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.save.New"
		log.WithFields(
			logrus.Fields{
				"op":         op,
				"request_id": requestid.GetReqID(c),
			}).Info("")

		//middleware.GetReqID(r.Context())

		var req Request
		err := c.ShouldBindJSON(&req)

		// err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// 	// Такую ошибку встретим, если получили запрос с пустым телом.
			// 	// Обработаем её отдельно
			log.Error("request body is empty")
			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("пустой запрос"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body.", err)
			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("невозможно декодировать запрос"))
			return
		}

		//		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("wrong request", err)
			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		_, err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", "url", req.URL)

			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error(storage.ErrURLExists.Error()))

			return
		}
		if err != nil {
			log.Error("failed to add url: ", err)

			c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("невозможно добавить ссылку"))

			return
		}

		log.Info("url added")

		// responseOK(w, r, alias)
		c.JSON(http.StatusOK, alias)

	}
}

// func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
// 	render.JSON(w, r, Response{
// 		Response: resp.OK(),
// 		Alias:    alias,
// 	})
// }

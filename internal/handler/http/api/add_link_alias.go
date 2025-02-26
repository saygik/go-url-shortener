package api

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	resp "github.com/saygik/go-url-shortener/internal/lib/api/response"
	"github.com/saygik/go-url-shortener/internal/lib/random"
	"github.com/saygik/go-url-shortener/internal/storage"
)

const aliasLength = 6

var (
	ErrURLNotFound = errors.New("адрес с таким псевдонимом не найден")
	ErrURLExists   = errors.New("адрес с таким псевдонимом уже существует")
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

func (h *Handler) SaveURL(c *gin.Context) {

	var req Request
	err := c.ShouldBindJSON(&req)
	if errors.Is(err, io.EOF) {
		// 	* Такую ошибку встретим, если получили запрос с пустым телом.
		// 	* Обработаем её отдельно
		h.log.Error("request body is empty")
		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("пустой запрос"))
		return
	}
	if err != nil {
		h.log.Error("failed to decode request body.", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("невозможно декодировать запрос"))
		return
	}
	req.URL = strings.ReplaceAll(req.URL, "+", "%20")
	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		h.log.Error("wrong request", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.ValidationError(validateErr))

		return
	}

	alias := req.Alias
	if alias == "" {
		alias = random.NewRandomString(aliasLength)
	}
	_, err = h.uc.SaveURL(req.URL, alias)
	if errors.Is(err, storage.ErrURLExists) {
		h.log.Info("url already exists", "url", req.URL)

		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error(storage.ErrURLExists.Error()))

		return
	}
	if err != nil {
		h.log.Error("failed to add url: ", err)

		c.AbortWithStatusJSON(http.StatusNotAcceptable, resp.Error("невозможно добавить ссылку"))

		return
	}

	h.log.Info("url added")

	// responseOK(w, r, alias)
	c.JSON(http.StatusOK, alias)
}

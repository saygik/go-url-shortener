package main

import (
	"github.com/saygik/go-url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg
}

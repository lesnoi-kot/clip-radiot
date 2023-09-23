package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func NewServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/cut", cutAudioHandler) // TODO: rate limit
	e.Static("/", "public")            // TODO: security

	return e
}

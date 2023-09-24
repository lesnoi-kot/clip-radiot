package server

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lesnoi-kot/clip-radiot/public"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func NewServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/cut", cutAudioHandler, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	if os.Getenv("DEBUG") != "" {
		e.Static("/", "public")
	} else {
		e.StaticFS("/", public.StaticData)
	}

	return e
}

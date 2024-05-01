package server

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lesnoi-kot/clip-radiot/public"
)

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

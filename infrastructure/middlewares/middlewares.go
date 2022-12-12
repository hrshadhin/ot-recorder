package middlewares

import (
	"ot-recorder/infrastructure/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const EchoLogFormat = "time: ${time_rfc3339_nano} || ${method}: ${uri} || u_agent: ${user_agent} || status: ${status}" +
	" || latency: ${latency_human} \n"

// Attach middlewares required for the application
func Attach(e *echo.Echo) error {
	cfg := config.Get().App

	// echo middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: EchoLogFormat}))
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit(cfg.RequestBodyLimit))
	e.Pre(middleware.RemoveTrailingSlash())

	return nil
}

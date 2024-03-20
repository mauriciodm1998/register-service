package rest

import (
	"register-service/internal/config"

	"github.com/labstack/echo/v4"
)

type Register interface {
	ClockIn(c echo.Context) error
}
type rest struct {
	register Register
}

func New(channel Register) rest {
	return rest{
		register: channel,
	}
}

func (r rest) Start() error {
	router := echo.New()

	mainGroup := router.Group("/api")
	registerGroup := mainGroup.Group("/clock-in")
	registerGroup.POST("/", r.register.ClockIn)
	//registerGroup.Use(middlewares.Authorization)

	return router.Start(":" + config.Get().Server.Port)
}

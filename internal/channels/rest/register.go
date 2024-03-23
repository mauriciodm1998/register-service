package rest

import (
	"context"
	"fmt"
	"net/http"
	"register-service/internal/channels"
	"register-service/internal/config"
	"register-service/internal/middlewares"
	"register-service/internal/service"
	"register-service/internal/token"

	"github.com/labstack/echo/v4"
)

type register struct {
	service service.RegisterService
}

func NewRegisterChannel() channels.Channel {
	return &register{
		service: service.NewRegisterService(),
	}
}

func (r *register) Start() error {
	router := echo.New()

	mainGroup := router.Group("/api")
	registerGroup := mainGroup.Group("/clock-in")

	mainGroup.GET("/healthz", r.healthCheck)
	registerGroup.POST("/", r.ClockIn)
	registerGroup.GET("/", r.GetDayAppointments)
	registerGroup.GET("/week", r.GetWeekAppointments)
	registerGroup.GET("/month", r.GetMonthAppointments)

	registerGroup.Use(middlewares.Authorization)
	registerGroup.Use(middlewares.Logger)

	return router.Start(":" + config.Get().Server.Port)
}

func (r *register) healthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (r *register) ClockIn(c echo.Context) error {
	userId, _, err := token.ExtractTokenValues(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid user").Error(),
		})
	}

	err = r.service.ClockIn(context.Background(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (r *register) GetDayAppointments(c echo.Context) error {
	userId, _, err := token.ExtractTokenValues(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid user").Error(),
		})
	}

	dailyRegisters, err := r.service.GetDayAppointments(context.Background(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, toResponse(*dailyRegisters))
}

func (r *register) GetWeekAppointments(c echo.Context) error {
	userId, _, err := token.ExtractTokenValues(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid user").Error(),
		})
	}

	weekRegisters, err := r.service.GetWeekAppointments(context.Background(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, toResponses(weekRegisters))
}

func (r *register) GetMonthAppointments(c echo.Context) error {
	userId, email, err := token.ExtractTokenValues(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid user").Error(),
		})
	}

	err = r.service.GetMonthAppointments(context.Background(), userId, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusCreated)
}

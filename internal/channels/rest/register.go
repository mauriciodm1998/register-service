package rest

import (
	"context"
	"fmt"
	"net/http"
	"register-service/internal/service"

	"github.com/labstack/echo/v4"
)

type register struct {
	service service.RegisterService
}

func NewRegisterChannel() Register {
	return &register{
		service: service.NewRegisterService(),
	}
}

func (r *register) ClockIn(c echo.Context) error {
	var registerRequest RegisterRequest

	if err := c.Bind(&registerRequest); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid data").Error(),
		})
	}

	err := r.service.ClockIn(context.Background(), registerRequest.ToClockInRegister())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

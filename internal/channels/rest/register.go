package rest

import (
	"context"
	"fmt"
	"net/http"
	"register-service/internal/service"
	"register-service/internal/token"

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
	userId, err := token.ExtractUserId(c.Request())
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

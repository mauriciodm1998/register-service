package middlewares

import (
	"context"
	"net/http"
	"register-service/internal/token"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Logger(fx echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Logger.WithContext(context.Background())
		request := ctx.Request()
		log.Info().Msgf("Host: %s, URI: %s, Method: %s", request.Host, request.RequestURI, request.Method)
		return fx(ctx)
	}
}

func Authorization(fx echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if err := token.ValidateToken(ctx.Request()); err != nil {
			ctx.Response().Header().Set("Content-Type", "application/json")
			ctx.Response().WriteHeader(http.StatusUnauthorized)
			return err
		}

		return fx(ctx)
	}
}

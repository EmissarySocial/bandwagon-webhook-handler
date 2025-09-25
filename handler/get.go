package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetPage() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusNoContent)
	}
}

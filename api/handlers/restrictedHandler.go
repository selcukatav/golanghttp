package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Restricted(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to protected area!! ")
}

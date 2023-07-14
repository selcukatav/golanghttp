package middlewares

import (
	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
)

func SetMainMiddleware(e *echo.Echo) {

	e.Use(serverHeader)

}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Selele/0.1")
		return next(c)
	}
}

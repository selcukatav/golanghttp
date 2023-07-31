package api

import (
	"goserver/api/handlers"
	"goserver/api/middlewares"

	"github.com/labstack/echo/v4"
)

func MainGroup(e *echo.Echo) {
	e.GET("/", handlers.MainPage)
}
func ResGroup(g *echo.Group) {
	g.GET("/res", middlewares.Authorize(handlers.Restricted))
}

package api

import (
	"goserver/api/handlers"
	"goserver/api/middlewares"

	"github.com/labstack/echo/v4"
)

func ResGroup(g *echo.Group) {
	g.GET("/res", middlewares.Authorize(handlers.Restricted))
}

package router

import (
	"goserver/api"
	"goserver/api/middlewares"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	//create groups
	adminGroup := e.Group("/admin")

	cookieGroup := e.Group("cookie")

	jwtGroup := e.Group("/jwt")

	//set all middlewares
	middlewares.SetMainMiddleware(e)
	middlewares.SetAdminMiddleware(adminGroup)
	middlewares.SetCookieMiddleware(cookieGroup)
	middlewares.SetJwtMiddleware(jwtGroup)

	//set main routes

	api.MainGroup(e)

	//set group routes

	api.AdminGroup(adminGroup)
	api.CookieGroup(cookieGroup)
	api.JwtGroup(jwtGroup)
	return e
}

package api

import (
	"goserver/api/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/", handlers.MainPage)
	e.GET("/users", handlers.GetUsers)
	e.GET("/login", handlers.Login)

	e.POST("/users", handlers.AddUser)
	e.POST("/admins", handlers.AddAdmin)
	e.POST("/moderators", handlers.AddMods)
}

package router

import (
	"goserver/api"
	"goserver/api/handlers"
	"goserver/database"

	"github.com/labstack/echo/v4"
)

func New(client *database.Mongo) *echo.Echo {
	e := echo.New()

	r := handlers.AuthRoute{
		Client: *client,
	}

	resGroup := e.Group("/restricted")

	e.GET("/", handlers.MainPage)

	e.GET("/login", handlers.GetLoginPage)
	e.GET("/signup", handlers.GetSignupPage)

	e.POST("/login", r.Login)
	e.POST("/signup", handlers.SignUpHandler(client))

	e.GET("/resetpassword", handlers.GetResetPassword)
	e.POST("/resetpassword", handlers.ResetPasswordHandler(client))

	api.ResGroup(resGroup)
	api.MainGroup(e)

	return e
}

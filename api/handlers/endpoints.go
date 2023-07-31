package handlers

import (
	"goserver/templates"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Restricted(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to protected area!! ")
}

func MainPage(c echo.Context) error {

	return c.String(http.StatusOK, "Welcome to Main Page")

}
func GetSignupPage(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "signup.html", nil)
}
func GetResetPassword(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "resetpassword.html", nil)
}
func GetLoginPage(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "login.html", nil)

}

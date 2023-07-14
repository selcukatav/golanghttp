package handlers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Main Page")

}
func MainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to Cookie Page")
}
func MainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to Admin Page")
}
func MainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	log.Println("User Name: ", claims["name"], "User ID", claims["jti"])

	return c.String(http.StatusOK, "You are reached the JWT")
}

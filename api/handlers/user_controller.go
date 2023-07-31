package handlers

import (
	"context"
	"goserver/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

func signUp(c echo.Context, client *database.Mongo) error {

	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")
	u := &User{

		Username: username,
		Password: password,
		Role:     role,
	}

	collection := client.Client.Database("GoServer").Collection("users")
	_, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return c.String(http.StatusBadRequest, "User Couldn't Add")
	}
	return c.String(http.StatusOK, "User successfully added! ")

}

func SignUpHandler(client *database.Mongo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return signUp(c, client)

	}
}

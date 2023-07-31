package handlers

import (
	"context"
	"fmt"
	"goserver/api/middlewares"
	"goserver/database"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`

	Role string `bson:"role"`
}

type AuthRoute struct {
	Client database.Mongo
}

func (a *AuthRoute) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	// Verifies user data if it's correct or not
	user, err := middlewares.Authenticate(&a.Client, username, password, role)
	if err != nil {
		return echo.ErrUnauthorized
	}

	// Generates both access token and refresh token
	accessToken, err := middlewares.GenerateToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	refreshToken, err := middlewares.GenerateRefreshToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	// Sets the access token as a response header
	c.Response().Header().Set("Authorization", "Bearer "+accessToken)

	// Creates a cookie for refresh token
	middlewares.SetCookieMiddleware(c, "access_token", accessToken, time.Now().Add(time.Hour*24*7))
	middlewares.SetCookieMiddleware(c, "refresh_token", refreshToken, time.Now().Add(time.Hour*24*7))

	// Return the access token in the JSON response
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func resetPassword(c echo.Context, client *database.Mongo) error {

	collection := client.Client.Database("GoServer").Collection("users")

	username := c.FormValue("username")
	newPassword := c.FormValue("password")
	againPassword := c.FormValue("newpassword")

	//filter gets the search parameter for mongodb
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"password": newPassword}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if newPassword == againPassword {
		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if result.ModifiedCount == 0 {
			return fmt.Errorf("User not found")
		}
		return c.String(http.StatusOK, "Your password successfuly updated.")
	}

	return nil

}

func ResetPasswordHandler(client *database.Mongo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return resetPassword(c, client)

	}
}

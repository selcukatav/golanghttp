package handlers

import (
	"context"
	"fmt"
	"goserver/database"
	"goserver/templates"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

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

func GetResetPassword(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "resetpassword.html", nil)
}

func ResetPasswordHandler(client *database.Mongo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return resetPassword(c, client)

	}
}

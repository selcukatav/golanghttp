package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Moderator struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}

func AddMods(c echo.Context) error {
	moderator := Moderator{}
	err := c.Bind(&moderator)
	if err != nil {
		log.Printf("Failed to add moderator at addMods")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("The User is: %#v ", moderator)
	return c.String(http.StatusOK, "The Moderator is added.")
}

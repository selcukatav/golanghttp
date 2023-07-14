package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Admin struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}

func AddAdmin(c echo.Context) error {
	admin := Admin{}

	err := json.NewDecoder(c.Request().Body).Decode(&admin)

	if err != nil {
		log.Printf("Failed to add admin at addAdmin")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer c.Request().Body.Close()
	log.Printf("The User is: %#v ", admin)
	return c.String(http.StatusOK, "The Admin is added.")
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}

func GetUsers(c echo.Context) error {
	userId := c.QueryParam("id")
	userName := c.QueryParam("name")
	userPwd := c.QueryParam("pwd")

	return c.String(http.StatusOK, fmt.Sprintf("Your Id : \n"+userId+"Your Name:\n "+userName+"Your password: \n"+userPwd))
}

func AddUser(c echo.Context) error {
	user := User{}

	b, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		log.Printf("Failed to reading the body in addUser ReadAll")
		return c.String(http.StatusInternalServerError, "")
	}
	defer c.Request().Body.Close()
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Printf("Failed to add user at addUser")
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("The User is: %#v ", user)
	return c.String(http.StatusOK, "The User is added.")
}

package server

//Get ve Post işlemleri için

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Id       string `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}

func GetUsers(c echo.Context) error {
	userId := c.QueryParam("id")
	userName := c.QueryParam("name")
	userPwd := c.QueryParam("pwd")

	return c.String(http.StatusOK, fmt.Sprintf("Your Id : %s, Your Name: %s, Your password: %s, ", userId, userName, userPwd))
}

func AddUser(c echo.Context) error {
	user := User{
		Id:       "1",
		UserName: "selcuk",
		UserPwd:  "12345",
	}

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

func MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to main page")
}

func Main() {
	e := echo.New()
	e.GET("/users", GetUsers)
	e.POST("/users", AddUser)
	e.GET("/", MainPage)

	e.Start(":8001")
}

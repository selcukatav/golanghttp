package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}
type Admin struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}
type Moderator struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	UserPwd  string `json:"pwd"`
}

var Users []User

func mainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Main Page")

}

func getUsers(c echo.Context) error {
	userId := c.QueryParam("id")
	userName := c.QueryParam("name")
	userPwd := c.QueryParam("pwd")

	return c.String(http.StatusOK, fmt.Sprintf("Your Id : \n"+userId+"Your Name:\n "+userName+"Your password: \n"+userPwd))
}

func addUser(c echo.Context) error {
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

func addAdmin(c echo.Context) error {
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

func addMods(c echo.Context) error {
	moderator := Moderator{}
	err := c.Bind(&moderator)
	if err != nil {
		log.Printf("Failed to add moderator at addMods")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("The User is: %#v ", moderator)
	return c.String(http.StatusOK, "The Moderator is added.")
}

func Demo1() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.GET("/", mainPage)
	e.GET("/users", getUsers)

	e.POST("/users", addUser)
	e.POST("/admins", addAdmin)
	e.POST("/moderators", addMods)

	e.Start(":8001")
}

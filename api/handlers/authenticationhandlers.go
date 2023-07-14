package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.Claims
}

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	//check username and password authentication from db

	if username == "selcuk" && password == "1234" {
		cookie := &http.Cookie{}

		//same thing
		//cookie:=new(http.Cookie)
		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)
		c.SetCookie(cookie)

		// TODO: Create jwt Token
		token, err := createJwtToken()
		if err != nil {
			log.Println("ERROR! Token Couldnt Created")
			return c.String(http.StatusInternalServerError, "Something went bad")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"massage": "You Logged In",
			"token":   token,
		})
	}
	return c.String(http.StatusUnauthorized, "Your username or password is invalid.")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"selcuk",
		jwt.MapClaims{
			"sub": "main_user_id",
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}
	return token, nil
}

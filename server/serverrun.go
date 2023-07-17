package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	userName string
	password string
	role     string
}

var users = []User{
	{userName: "selcuk", password: "12345", role: "admin"},
}

/*type jwtCustomClaims struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}*/

var jwtKey interface{} = []byte("my_secret_key")

func Login(c echo.Context) error {
	username := "selcuk"
	password := "12345"
	role := "admin"

	user, err := authenticate(username, password, role)
	if err != nil {
		return echo.ErrUnauthorized
	}

	tokenString, err := generateToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Set("user", tokenString)

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 4),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})

}

func getJWTKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("error occurred at authorization")
	}
	return jwtKey, nil
}

func authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.ErrUnauthorized
		}

		tokenString := strings.Split(authHeader, "Bearer")[1]
		cleanedToken := strings.TrimSpace(tokenString)

		token, err := jwt.Parse(cleanedToken, getJWTKey)
		if err != nil {
			return echo.ErrUnauthorized
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role := claims["role"].(string)
			if role != "admin" {
				return echo.ErrForbidden
			}
			return next(c)
		}
		return echo.ErrUnauthorized
	}
}

func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.userName,
		"role":     user.role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func authenticate(username, password, role string) (*User, error) {

	for _, user := range users {
		if user.userName == username && user.password == password && user.role == role {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("error occured")

}
func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Yes it is accessible")
}

func restricted(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to protected area!! ")
}

func Main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", Login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")
	//r.Use(middleware.JWT([]byte(jwtKey)))

	r.GET("/res", authorize(restricted))

	e.Logger.Fatal(e.Start(":1323"))
}

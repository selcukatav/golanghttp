package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`

	Role string `bson:"role"`
}

var jwtKey = []byte("my-secret-key")

func GenerateToken(user *User) (string, error) {
	claims := jwt.MapClaims{

		"sub":      1,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return accessToken, nil

}

func GenerateRefreshToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      1,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 4).Unix(), // Örneğin, 1 haftalık bir süre belirleyebilirsiniz.
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func SetCookieMiddleware(c echo.Context, name, value string, expires time.Time) {
	cookie := &http.Cookie{

		Name:    name,
		Value:   value,
		Expires: expires,
		Path:    "/",
	}

	c.SetCookie(cookie)
}

package middlewares

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func Authorize(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		//Take the authorization request key value
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			tokenHeader := authHeader[len("Bearer "):]

			token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

					return nil, echo.ErrUnauthorized
				}
				return jwtKey, nil
			})
			if err == nil && token.Valid {
				return c.String(http.StatusOK, "You reached the location successfully")
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				role := claims["role"].(string)
				if role != "admin" {
					return echo.ErrForbidden
				}
				return next(c)
			}
		}

		cookie, err := c.Cookie("token")
		if err == nil {
			tokenCookie := cookie.Value

			token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error occurred at authorization")
				}
				return jwtKey, nil
			})
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
		}

		return echo.ErrUnauthorized
	}
}

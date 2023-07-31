package middlewares

import (
	"context"
	"fmt"
	"goserver/database"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Authenticate(client *database.Mongo, username, password, role string) (*User, error) {

	if client == nil {
		fmt.Println("client değişkeni nil.")

	}

	collection := client.Client.Database("GoServer").Collection("users")

	filter := bson.M{"username": username, "password": password, "role": role}

	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("username or password is worng")
		}
		return nil, fmt.Errorf("error occured in authenticate")
	}
	return &user, nil

}

func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Take the authorization request key value
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
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					sub := int(claims["sub"].(float64))

					if sub != 1 {
						return echo.ErrForbidden
					}
					return next(c)

				}
			}
		}

		cookie, _ := c.Cookie("access_token")
		if cookie != nil {
			tokenCookie := cookie.Value

			token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error occurred at authorization")
				}
				return jwtKey, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					role := claims["role"].(string)
					if role != "admin" {
						return echo.ErrForbidden
					}
					return next(c)
				}
			}
		}

		return echo.ErrUnauthorized
	}
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
}

var jwtKey = []byte("my_secret_key")

var client *mongo.Client

func Init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("MongoDb connection error", err)
		return
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		fmt.Println("Mongodb Ping Error", err)
		return
	}

	fmt.Println("Successfully connected to the MongoDb.")
}
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	//verifies user data if its correct or not

	user, err := authenticate(username, password, role)
	if err != nil {
		return echo.ErrUnauthorized
	}

	//generates token

	tokenString, err := generateToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Set("user", tokenString)

	//cookie elements
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 4),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	//creates cookie
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		//outputs the tokenString
		"token": tokenString,
	})

}

func authorize(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		//Take the authorization request key value
		cookie, err := c.Cookie("token")
		if err != nil {
			return echo.ErrUnauthorized
		}
		//
		tokenString := cookie.Value
		cleanedToken := strings.TrimSpace(tokenString)

		token, err := jwt.Parse(cleanedToken, func(token *jwt.Token) (interface{}, error) {
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
		return echo.ErrUnauthorized
	}
}

// create new token func
func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//authenticate the user

func authenticate(username, password, role string) (*User, error) {

	if client == nil {
		fmt.Println("client değişkeni nil.")

	}

	collection := client.Database("GoServer").Collection("users")

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
func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Yes it is accessible")
}

func restricted(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to protected area!! ")
}

func Main() {
	e := echo.New()

	Init()

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

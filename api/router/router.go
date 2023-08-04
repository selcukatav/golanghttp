package router

import (
	"context"
	"fmt"
	"goserver/templates"
	"net/http"
	"time"

	"goserver/api/middlewares"
	"goserver/database"
	_ "goserver/docs"

	"github.com/go-gomail/gomail"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`

	Role string `bson:"role"`
}

type AuthRoute struct {
	Client database.Mongo
}

// @Summary Kullanıcı kaydı
// @Description Yeni bir kullanıcı hesabı oluşturma işlemi
// @ID user-signup
// @Accept json
// @Produce json
// @Param user body User true "Kullanıcı bilgileri"
// @Success 201 {string} string "Kullanıcı başarıyla oluşturuldu!"
// @Failure 400 {string} string "Hatalı istek"
// @Router /signup [post]
func signUp(c echo.Context, client *database.Mongo) error {

	email := c.FormValue(("email"))

	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")
	u := &User{
		Email:    email,
		Username: username,
		Password: password,
		Role:     role,
	}

	collection := client.Client.Database("GoServer").Collection("users")
	_, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return c.String(http.StatusBadRequest, "User Couldn't Add")
	}
	return c.String(http.StatusOK, "User successfully added! ")

}

func SignUpHandler(client *database.Mongo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return signUp(c, client)

	}
}

// @Summary Kullanıcı girişi
// @Description Kullanıcı adı ve şifre ile giriş yapma işlemi
// @ID user-login
// @Accept json
// @Produce json
// @Param username query string true "Kullanıcı adı"
// @Param password query string true "Şifre"
// @Param role query string true "Rol"
// @Success 200 {object} map[string]string "Giriş başarılı!"
// @Failure 401 {object} string "Yetkilendirme hatası"
// @Router /login [post]
func (a *AuthRoute) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	// Verifies user data if it's correct or not
	user, err := middlewares.Authenticate(&a.Client, username, password, role)
	if err != nil {
		return echo.ErrUnauthorized
	}

	// Generates both access token and refresh token
	accessToken, err := middlewares.GenerateToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	refreshToken, err := middlewares.GenerateRefreshToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	// Sets the access token as a response header
	c.Response().Header().Set("Authorization", "Bearer "+accessToken)

	// Creates a cookie for refresh token
	middlewares.SetCookieMiddleware(c, "access_token", accessToken, time.Now().Add(time.Hour*24*7))
	middlewares.SetCookieMiddleware(c, "refresh_token", refreshToken, time.Now().Add(time.Hour*24*7))

	// Return the access token in the JSON response
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary Şifre sıfırlama
// @Description Kullanıcının şifresini sıfırlama işlemi
// @ID user-reset-password
// @Accept json
// @Produce json
// @Param username query string true "Kullanıcı adı"
// @Param password query string true "Yeni şifre"
// @Param newpassword query string true "Yeni şifre (tekrar)"
// @Success 200 {string} string "Şifre başarıyla güncellendi!"
// @Failure 400 {object} string "Hatalı istek"
// @Router /resetpassword [post]
func resetPassword(c echo.Context, client *database.Mongo) error {

	collection := client.Client.Database("GoServer").Collection("users")

	username := c.FormValue("username")
	newPassword := c.FormValue("password")
	againPassword := c.FormValue("newpassword")

	//filter gets the search parameter for mongodb
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"password": newPassword}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if newPassword == againPassword {
		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if result.ModifiedCount == 0 {
			return fmt.Errorf("User not found")
		}
		return c.String(http.StatusOK, "Your password successfuly updated.")
	}

	return nil

}

func sendEmail(c echo.Context, client *database.Mongo) error {
	email := c.FormValue("email")

	collection := client.Client.Database("GoServer").Collection("users")
	filter := bson.D{{"email", email}}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return c.String(http.StatusNotFound, "Could not find the e-mail address")
	} else if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "selcukatav41@hotmail.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Reset Password")
	mailer.SetBody("text/html", "Click the link below to reset your password <a href=\"http://localhost:1323/resetpassword</a>")
	dialer := gomail.NewDialer("mail.example.com", 587, "selcukatav41@hotmail.com", "bd4e7132")
	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println("E-posta gönderme hatası:", err)
		return c.String(http.StatusInternalServerError, "E-posta gönderimi başarısız")
	}

	fmt.Println("E-posta gönderildi:", email)
	return c.String(http.StatusOK, "E-posta gönderildi: "+email)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func New(client *database.Mongo) *echo.Echo {
	e := echo.New()

	r := AuthRoute{
		Client: *client,
	}

	resGroup := e.Group("/restricted")

	e.GET("/", MainPage)

	e.GET("/login", GetLoginPage)
	e.GET("/signup", GetSignupPage)

	e.POST("/login", r.Login)

	e.POST("/signup", SignUpHandler(client))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/resetpassword", GetResetPassword)
	e.POST("/resetpassword", ResetPasswordHandler(client))

	e.GET("/sendemail", SendEmailHandler(client))
	e.POST("/sendemail", GetResetPasswordEmailPage)

	ResGroup(resGroup)
	MainGroup(e)

	return e
}

func ResetPasswordHandler(client *database.Mongo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return resetPassword(c, client)

	}
}
func SendEmailHandler(client *database.Mongo) echo.HandlerFunc {
	return func(c echo.Context) error {
		return sendEmail(c, client)

	}
}
func Restricted(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to protected area!! ")
}

// @Summary Greet the user
// @Description Say hello to the user
// @ID hello-world
// @Produce json
// @Success 200 {object} string "Hello, World!"
// @Router / [get]
func MainPage(c echo.Context) error {

	return c.String(http.StatusOK, "Welcome to Main Page")

}

func GetResetPasswordEmailPage(c echo.Context) error {
	return templates.Templates.ExecuteTemplate(c.Response(), "resetpasswordemail.html", nil)
}
func GetSignupPage(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "signup.html", nil)
}
func GetResetPassword(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "resetpassword.html", nil)
}
func GetLoginPage(c echo.Context) error {

	return templates.Templates.ExecuteTemplate(c.Response(), "login.html", nil)

}
func MainGroup(e *echo.Echo) {
	e.GET("/", MainPage)
}
func ResGroup(g *echo.Group) {
	g.GET("/res", middlewares.Authorize(Restricted))
}

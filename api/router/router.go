package router

import (
	"goserver/api"
	"goserver/api/handlers"
	"goserver/database"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

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

	r := handlers.AuthRoute{
		Client: *client,
	}

	resGroup := e.Group("/restricted")

	// @Summary Greet the user
	// @Description Say hello to the user
	// @ID hello-world
	// @Produce json
	// @Success 200 {object} string "Hello, World!"
	// @Router / [get]

	e.GET("/", handlers.MainPage)

	e.GET("/login", handlers.GetLoginPage)
	e.GET("/signup", handlers.GetSignupPage)

	// @Summary Kullanıcı girişi
	// @Description Kullanıcı adı ve şifre ile giriş yapma işlemi
	// @ID user-login
	// @Accept json
	// @Produce json
	// @Param credentials body LoginCredentials true "Giriş bilgileri"
	// @Success 200 {object} string "Giriş başarılı!"
	// @Failure 400 {object} string "Hatalı istek"
	// @Router /login [post]
	e.POST("/login", r.Login)

	// @Summary Kullanıcı kaydı
	// @Description Yeni bir kullanıcı hesabı oluşturma işlemi
	// @ID user-signup
	// @Accept json
	// @Produce json
	// @Param user body SignupUser true "Kullanıcı bilgileri"
	// @Success 201 {object} string "Kullanıcı başarıyla oluşturuldu!"
	// @Failure 400 {object} string "Hatalı istek"
	// @Router /signup [post]

	e.POST("/signup", handlers.SignUpHandler(client))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/resetpassword", handlers.GetResetPassword)
	e.POST("/resetpassword", handlers.ResetPasswordHandler(client))

	api.ResGroup(resGroup)
	api.MainGroup(e)

	return e
}

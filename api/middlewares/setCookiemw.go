package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetCookieMiddleware(c echo.Context, name, value string, expires time.Time) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
		Path:    "/",
	}

	c.SetCookie(cookie)
}

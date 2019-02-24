package handler

import (
	"net/http"
	"time"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "test" && password == "test" {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "test"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func Auth() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)

		return c.String(http.StatusOK, "Welcome " + name + "!")
	}
}
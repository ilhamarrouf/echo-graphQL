package handler

import (
	"bytes"
	"log"
	"net/http"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"github.com/ilhamarrouf/echo-graphql/graphql"
	"github.com/ilhamarrouf/echo-graphql/middlewares"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	tokenString, refreshTokenString, err := middlewares.CreateNewTokens(username, password)
	if err == nil {
		return c.JSON(http.StatusOK, map[string]string{
			"token": tokenString,
			"refresh_token": refreshTokenString,
		})
	}

	return echo.ErrUnauthorized
}

func Restricted() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		_ = user.Claims.(jwt.MapClaims)
		bufBody := new(bytes.Buffer)
		bufBody.ReadFrom(c.Request().Body)
		query := bufBody.String()
		log.Printf(query)
		result := graphql.ExecuteQuery(query)

		return c.JSON(http.StatusOK, result)
	}
}

func ReAuth() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*middlewares.MyClaim)
		oldToken := c.FormValue("old_token")
		tokenString, refreshTokenString, err := middlewares.UpdateRefreshTokenExp(claims, oldToken)
		if err == nil {
			return c.JSON(http.StatusOK, map[string]string{
				"token": tokenString,
				"refresh_token": refreshTokenString,
			})
		}

		return echo.ErrUnauthorized
	}
}
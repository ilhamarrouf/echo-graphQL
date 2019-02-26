package handler

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"github.com/ilhamarrouf/echo-graphql/db"
	"github.com/ilhamarrouf/echo-graphql/models"
	"github.com/ilhamarrouf/echo-graphql/graphql"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	db := db.CreateConnection()
	db.SingularTable(true)
	user := [] models.User{}
	db.Find(&user, "name=? and password=?", username, password)

	if len(user) > 0 && username == user[0].Name {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
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
		_ = user.Claims.(jwt.MapClaims)
		bufBody := new(bytes.Buffer)
		bufBody.ReadFrom(c.Request().Body)
		query := bufBody.String()
		log.Printf(query)
		result := graphql.ExecuteQuery(query)

		return c.JSON(http.StatusOK, result)
	}
}
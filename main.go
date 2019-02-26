package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ilhamarrouf/echo-graphql/handler"
	"net/http"
)

func main()  {
	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE},
	}))

	app.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Echo Framework</h1>")
	})
	app.POST("/login", handler.Login)

	auth := app.Group("/auth")
	auth.Use(middleware.JWT([]byte("secret")))
	auth.POST("", handler.Auth())

	app.Logger.Fatal(app.Start(":3000"))
}
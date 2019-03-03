package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ilhamarrouf/echo-graphql/handler"
	ext "github.com/ilhamarrouf/echo-graphql/middlewares"
	"github.com/ilhamarrouf/echo-graphql/libs"
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

	libs.InitDB()

	app.GET("/hello", handler.Hello())
	app.POST("/login", handler.Login)
	r1 := app.Group("/restricted")
	r1.Use(middleware.JWT([]byte("secret1")))
	r1.POST("", handler.Restricted())

	r2 := app.Group("/reauth")
	config := middleware.JWTConfig{
		Claims:     &ext.MyClaim{},
		SigningKey: []byte("secret2"),
	}
	r2.Use(middleware.JWTWithConfig(config))
	r2.POST("", handler.ReAuth())

	app.Logger.Fatal(app.Start(":3000"))
}
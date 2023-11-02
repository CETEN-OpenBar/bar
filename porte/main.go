package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// CORS
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Local-Token")
			return next(c)
		}
	})

	routes(e)

	e.Logger.Fatal(e.Start(":8080"))
}

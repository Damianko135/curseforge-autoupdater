package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		// As per RFC 2324, a teapot responds with status code 418 I'm a Teapot
		return c.String(http.StatusTeapot, "I'm a teapot")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, CurseForge AutoUpdate!")
	})

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}

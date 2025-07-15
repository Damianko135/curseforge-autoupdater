package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/damianko135/curseforge-autoupdate/golang/views"
)

func main() {
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files
	e.Static("/static", "public")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return render(c, views.home())
	})

	e.GET("/health", func(c echo.Context) error {
		return render(c, views.health())
	})

	e.GET("/status", func(c echo.Context) error {
		return render(c, views.status())
	})

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}

// render is a helper function to render templ components
func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response().Writer)
}

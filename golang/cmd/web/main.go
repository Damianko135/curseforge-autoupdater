package main

import (
	"github.com/a-h/templ"
	"github.com/damianko135/curseforge-autoupdate/golang/views"  //nolint:all
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Add middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files
	e.Static("/static", "public")

	// Routes
	// NOTE: It will through an error if templ hasnt build the files yet.
	e.GET("/", func(c echo.Context) error {
		return render(c, views.Status())
	})

	e.GET("/health", func(c echo.Context) error {
		return render(c, views.Health())
	})

	e.GET("/status", func(c echo.Context) error {
		return render(c, views.Status())
	})

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}

// render is a helper function to render templ components
func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response().Writer)
}

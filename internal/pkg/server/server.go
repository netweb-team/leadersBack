package server

import (
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/database"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cfg := config.New()
	db := database.Connect(cfg.Postgres)
	e.Logger.Debug(db)

	// Routes
	e.GET("/api", hello)

	// Start server
	e.Logger.Fatal(e.Start(cfg.Port))
}

func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "hello")
}

package server

import (
	"crypto/tls"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/htmlparser"
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
	api := e.Group("/api")
	api.GET("", hello)
	api.GET("/html", parseHTML)
	api.GET("/html/ad", parseAdHTML)

	// Start server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	e.Logger.Fatal(e.Start(cfg.Port))
}

func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "hello")
}

func parseHTML(ctx echo.Context) error {
	result := htmlparser.Search(ctx.QueryParam("url"))
	return ctx.JSON(http.StatusOK, result)
}

func parseAdHTML(ctx echo.Context) error {
	result := htmlparser.Ad(ctx.QueryParam("url"), ctx.QueryParam("floor"))
	return ctx.JSON(http.StatusOK, result)
}

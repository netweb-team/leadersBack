package server

import (
	"crypto/tls"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/handler"
	"leaders_apartments/internal/pkg/parser/html"
	"leaders_apartments/internal/pkg/repository"
	"leaders_apartments/internal/pkg/usecase"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Run() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cfg := config.New()
	db := database.Connect(cfg.Postgres)
	repo := repository.New(db)
	uc := usecase.New(repo)
	handlers := handler.New(uc)

	// Routes
	api := e.Group("/api")
	api.GET("", hello)
	api.POST("/pools", handlers.ImportXslx)
	api.GET("/pools/:id", handlers.GetPool)

	// Start server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	e.Logger.Fatal(e.Start(cfg.Port))
}

func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "hello")
}

func parseHTML(ctx echo.Context) error {
	url := ctx.QueryParam("url")
	log.Info(url)
	result := html.Search(strings.ReplaceAll(url, ",", "&"))
	return ctx.JSON(http.StatusOK, result)
}

func parseAdHTML(ctx echo.Context) error {
	result := html.Ad(ctx.QueryParam("url"))
	return ctx.JSON(http.StatusOK, result)
}

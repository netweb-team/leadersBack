package server

import (
	"crypto/tls"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/handler"
	"leaders_apartments/internal/pkg/repository"
	"leaders_apartments/internal/pkg/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())

	cfg := config.New()
	db := database.Connect(cfg.Postgres)
	repo := repository.New(db)
	uc := usecase.New(repo)
	handlers := handler.New(uc)

	// Routes
	api := e.Group("/api")
	api.GET("", hello)
	api.POST("/signup", handlers.SignUp)
	api.POST("/auth", handlers.SignIn)
	api.DELETE("/auth", handlers.SignOut)
	api.POST("/pools", handlers.ImportXlsx)
	api.GET("/pools/:id", handlers.GetPool)
	api.POST("/pools/:id", handlers.CalcPool)
	api.PUT("/pools/:id", handlers.ChangePool)

	// Start server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	e.Logger.Fatal(e.Start(cfg.Port))
}

func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "hello")
}

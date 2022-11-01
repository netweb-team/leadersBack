package handler

import (
	"leaders_apartments/internal/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type serverHandlers struct {
	uc domain.Usecase
}

func New(uc domain.Usecase) domain.Handler {
	return &serverHandlers{uc}
}

func (h *serverHandlers) ImportXslx(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "")
}

package handler

import (
	"leaders_apartments/internal/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type serverHandlers struct {
	uc domain.Usecase
}

func New(uc domain.Usecase) domain.Handler {
	return &serverHandlers{uc}
}

func EchoResponse(ctx echo.Context, status int, body interface{}) error {
	return ctx.JSON(http.StatusOK, &domain.Response{Status: status, Body: body})
}

func (h *serverHandlers) ImportXslx(ctx echo.Context) error {
	formFile, err := ctx.FormFile("table")
	if err != nil {
		log.Error("No file in multipart ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	file, err := formFile.Open()
	if err != nil {
		log.Error("Cannot open file from form ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	defer file.Close()

	if result := h.uc.ImportXslx(file); result != nil {
		return EchoResponse(ctx, http.StatusCreated, result)
	}
	return EchoResponse(ctx, http.StatusInternalServerError, nil)
}

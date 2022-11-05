package handler

import (
	"leaders_apartments/internal/pkg/domain"
	"net/http"
	"strconv"

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

func (h *serverHandlers) ImportXlsx(ctx echo.Context) error {
	formFile, err := ctx.FormFile("table")
	if err != nil {
		log.Info("No file in multipart ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	file, err := formFile.Open()
	if err != nil {
		log.Info("Cannot open file from form ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	defer file.Close()

	if result := h.uc.ImportXlsx(file); result != nil {
		return EchoResponse(ctx, http.StatusCreated, result)
	}
	return EchoResponse(ctx, http.StatusInternalServerError, nil)
}

func (h *serverHandlers) GetPool(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Info("Bad id in url: ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	if ptn := ctx.QueryParam("pattern"); len(ptn) > 0 {
		if p, err := strconv.Atoi(ptn); err == nil {
			result := h.uc.FindAnalogs(id, p)
			if result == nil {
				return EchoResponse(ctx, http.StatusNotFound, nil)
			}
			return EchoResponse(ctx, http.StatusOK, result)
		}
		log.Info("Bad pattern query param: ", err)
	} else if ctx.QueryParam("download") == "1" {
		file := h.uc.ExportXlsx(id)
		if file == "" {
			return EchoResponse(ctx, http.StatusNotFound, nil)
		}
		return ctx.Attachment(file, file)
	}
	return EchoResponse(ctx, http.StatusBadRequest, nil)
}

func (h *serverHandlers) CalcPool(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Info("Bad id in url: ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	result := h.uc.CalcPool(id)
	if result == nil {
		return EchoResponse(ctx, http.StatusNotFound, nil)
	}
	return EchoResponse(ctx, http.StatusOK, result)
}

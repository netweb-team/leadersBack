package handler

import (
	"leaders_apartments/internal/pkg/domain"
	"net/http"
	"strconv"
	"time"

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
	user := 0
	if c, err := ctx.Cookie("session"); err == nil {
		user = h.uc.CheckAuth(c.Value, 0)
	}
	if user == 0 {
		return EchoResponse(ctx, http.StatusForbidden, nil)
	}
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

	if result := h.uc.ImportXlsx(file, user); result != nil {
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
	if c, err := ctx.Cookie("session"); err != nil || h.uc.CheckAuth(c.Value, id) == 0 {
		return EchoResponse(ctx, http.StatusForbidden, nil)
	}
	if ptn := ctx.QueryParam("pattern"); len(ptn) > 0 {
		if p, err := strconv.Atoi(ptn); err == nil {
			if result := h.uc.FindAnalogs(id, p); result != nil {
				return EchoResponse(ctx, http.StatusOK, result)
			}
		}
		log.Info("Bad pattern query param: ", err)
	} else if ctx.QueryParam("download") == "1" {
		if file := h.uc.ExportXlsx(id); file != "" {
			ctx.Response().Header().Set("Content-Description", "File Transfer")
			ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			return ctx.Attachment(file, file)
		}
	} else if result := h.uc.GetArchive(id); result != nil {
		return EchoResponse(ctx, http.StatusOK, result)
	}
	return EchoResponse(ctx, http.StatusNotFound, nil)
}

func (h *serverHandlers) CalcPool(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Info("Bad id in url: ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	if c, err := ctx.Cookie("session"); err != nil || h.uc.CheckAuth(c.Value, id) == 0 {
		return EchoResponse(ctx, http.StatusForbidden, nil)
	}
	result := h.uc.CalcPool(id)
	if result == nil {
		return EchoResponse(ctx, http.StatusNotFound, nil)
	}
	return EchoResponse(ctx, http.StatusOK, result)
}

func newCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    value,
		Expires:  time.Now().Add(240 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

func (h *serverHandlers) SignUp(ctx echo.Context) error {
	user := new(domain.User)
	if err := ctx.Bind(user); err != nil {
		log.Error(err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	value := h.uc.CreateUser(user)
	if value == "" {
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	ctx.SetCookie(newCookie(value))
	return EchoResponse(ctx, http.StatusCreated, nil)
}

func (h *serverHandlers) SignOK(ctx echo.Context) error {
	if c, err := ctx.Cookie("session"); err != nil || h.uc.CheckAuth(c.Value, 0) == 0 {
		return EchoResponse(ctx, http.StatusForbidden, nil)
	}
	return EchoResponse(ctx, http.StatusOK, nil)
}

func (h *serverHandlers) SignIn(ctx echo.Context) error {
	user := new(domain.User)
	if err := ctx.Bind(user); err != nil {
		log.Error(err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	value := h.uc.CreateAuth(user)
	if value == "" {
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	ctx.SetCookie(newCookie(value))
	return EchoResponse(ctx, http.StatusCreated, nil)
}

func (h *serverHandlers) SignOut(ctx echo.Context) error {
	if cookie, err := ctx.Cookie("session"); err == nil {
		h.uc.DeleteAuth(cookie.Value)
	}
	return EchoResponse(ctx, http.StatusOK, nil)
}

func (h *serverHandlers) ChangePool(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Info("Bad id in url: ", err)
		return EchoResponse(ctx, http.StatusBadRequest, nil)
	}
	if c, err := ctx.Cookie("session"); err != nil || h.uc.CheckAuth(c.Value, id) == 0 {
		return EchoResponse(ctx, http.StatusForbidden, nil)
	}
	if anal := ctx.QueryParam("correct"); len(anal) > 0 {
		if a, err := strconv.Atoi(anal); err == nil {
			coefs := new(domain.CorrectCoefs)
			if err = ctx.Bind(coefs); err != nil {
				log.Error(err)
				return EchoResponse(ctx, http.StatusBadRequest, nil)
			}
			if result := h.uc.ChangeCorrect(id, a, coefs); result != nil {
				return EchoResponse(ctx, http.StatusOK, result)
			}
		}
		log.Info("Bad pattern query param: ", err)
	} else if anal := ctx.QueryParam("analog"); len(anal) > 0 {
		if a, err := strconv.Atoi(anal); err == nil {
			if result := h.uc.ChangeAnalog(id, a); result != nil {
				return EchoResponse(ctx, http.StatusOK, result)
			}
		}
		log.Info("Bad pattern query param: ", err)
	}
	return EchoResponse(ctx, http.StatusNotFound, nil)
}

func (h *serverHandlers) GetUserArchives(ctx echo.Context) error {
	user := 0
	if c, err := ctx.Cookie("session"); err == nil {
		user = h.uc.CheckAuth(c.Value, 0)
	}
	if user == 0 {
		return EchoResponse(ctx, http.StatusForbidden, nil)
	}
	if result := h.uc.GetUserArchives(user); result != nil {
		return EchoResponse(ctx, http.StatusOK, result)
	}
	return EchoResponse(ctx, http.StatusNotFound, nil)
}

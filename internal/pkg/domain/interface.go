package domain

import (
	"os"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	SaveTable(filename string) error
}

type Usecase interface {
	ImportXslx(f *os.File) error
}

type Handler interface {
	ImportXslx(ctx echo.Context) error
}

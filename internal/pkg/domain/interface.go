package domain

import (
	"io"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	SaveTable(filename string) error
}

type Usecase interface {
	ImportXslx(f io.Reader) *Table
}

type Handler interface {
	ImportXslx(ctx echo.Context) error
}

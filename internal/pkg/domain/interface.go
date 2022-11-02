package domain

import (
	"io"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	SaveTable(filename string) (int, error)
}

type Usecase interface {
	ImportXslx(f io.Reader) *Table
	FindAnalogs(id, ptnIndex int) []*AdPage
}

type Handler interface {
	ImportXslx(ctx echo.Context) error
	GetPool(ctx echo.Context) error
}

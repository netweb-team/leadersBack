package domain

import (
	"io"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	SaveTable(filename string) (int, error)
	GetTableName(id int) (string, error)
}

type Usecase interface {
	ImportXslx(f io.Reader) *Table
	FindAnalogs(id, ptnIndex int) *PatternAnalogs
}

type Handler interface {
	ImportXslx(ctx echo.Context) error
	GetPool(ctx echo.Context) error
}

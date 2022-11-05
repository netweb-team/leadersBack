package domain

import (
	"io"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	SaveTable(filename string) (int, error)
	GetTableName(id int) (string, error)
	SaveAnalogs(id int, pattern *Row, analogs []*Row, coefs []*CorrectCoefs) error
	GetPatternAnalogs(id int) ([]*PatternAnalogs, error)
	GetAnalogs(id, ptnID int) ([]*Row, []*CorrectCoefs)
}

type Usecase interface {
	ImportXlsx(f io.Reader) *Table
	FindAnalogs(id, ptnIndex int) *PatternAnalogs
	CalcPool(id int) []*Row
	ExportXlsx(id int) string
}

type Handler interface {
	ImportXlsx(ctx echo.Context) error
	GetPool(ctx echo.Context) error
	CalcPool(ctx echo.Context) error
}

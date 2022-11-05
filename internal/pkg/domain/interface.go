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
	CreateCookie(id int, cookie string) error
	DeleteCookie(cookie string)
	CheckCookie(cookie string) int
	CreateUser(user *User) error
	GetUser(login string) *User
}

type Usecase interface {
	ImportXlsx(f io.Reader) *Table
	FindAnalogs(id, ptnIndex int) *PatternAnalogs
	CalcPool(id int) []*Row
	ExportXlsx(id int) string
	CreateUser(user *User) string
	CreateAuth(user *User) string
	DeleteAuth(cookie string)
	CheckAuth(cookie string) int
}

type Handler interface {
	ImportXlsx(ctx echo.Context) error
	GetPool(ctx echo.Context) error
	CalcPool(ctx echo.Context) error
	SignUp(ctx echo.Context) error
	SignIn(ctx echo.Context) error
	SignOut(ctx echo.Context) error
}

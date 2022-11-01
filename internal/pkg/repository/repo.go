package repository

import (
	"context"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/domain"

	"github.com/labstack/gommon/log"
)

const (
	insertTable = `insert into tables (path) values($1);`
)

type dbRepository struct {
	db *database.DBManager
}

func New(db *database.DBManager) domain.Repository {
	return &dbRepository{db}
}

func (repo *dbRepository) SaveTable(filename string) error {
	_, err := repo.db.Pool.Exec(context.Background(), insertTable, filename)
	if err != nil {
		log.Error("Unable to save path to table: ", err)
	}
	return err
}

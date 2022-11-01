package repository

import (
	"context"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/domain"

	"github.com/labstack/gommon/log"
)

const (
	insertTable = `insert into tables (path) values($1) returning id;`
)

type dbRepository struct {
	db *database.DBManager
}

func New(db *database.DBManager) domain.Repository {
	return &dbRepository{db}
}

func (repo *dbRepository) SaveTable(filename string) (int, error) {
	var id int
	err := repo.db.Pool.QueryRow(context.Background(), insertTable, filename).Scan(&id)
	if err != nil {
		log.Error("Unable to save path to table: ", err)
	}
	return id, err
}

package repository

import (
	"context"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/domain"

	"github.com/labstack/gommon/log"
)

const (
	insertTable     = `insert into tables (path) values($1) returning id;`
	selectTableName = `select path from tables where id = $1;`
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

func (repo *dbRepository) GetTableName(id int) (string, error) {
	var name string
	err := repo.db.Pool.QueryRow(context.Background(), selectTableName, id).Scan(&name)
	if err != nil {
		log.Error("Unable to get path to table: ", err)
	}
	return name, err
}

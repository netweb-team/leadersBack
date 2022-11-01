package repository

import (
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/domain"
)

type dbRepository struct {
	db *database.DBManager
}

func New(db *database.DBManager) domain.Repository {
	return &dbRepository{db}
}

func (repo *dbRepository) SaveTable(filename string) error {
	return nil
}

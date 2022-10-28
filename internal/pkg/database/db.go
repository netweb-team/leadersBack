package database

import (
	"context"
	"fmt"
	"leaders_apartments/internal/pkg/config"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type DBInterface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type DBManager struct {
	Pool DBInterface
}

func Connect(cfg *config.DBConfig) *DBManager {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Error("Connection to postgres failed")
		return nil
	}
	log.Info("Successful connection to postgres")
	return &DBManager{Pool: pool}
}

func Disconnect(db *DBManager) {
	db.Pool.Close()
	log.Info("database disconnected")
}

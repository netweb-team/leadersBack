package repository

import (
	"context"
	"fmt"
	"leaders_apartments/internal/pkg/database"
	"leaders_apartments/internal/pkg/domain"
	"sync"

	"github.com/labstack/gommon/log"
)

const (
	insertTable     = `insert into tables(path) values($1) returning id;`
	selectTableName = `select path from tables where id = $1;`
	insertPattern   = `insert into patterns(pool_id, pattern, lng, lat, avg_price) values($1, $2, $3, $4, $5);`
	insertAnalog    = `insert into analogs(lng,lat,addr,room,segment,floors,cur_floor,walls,total,kitchen,balcony,metro,state,price,avg_price,
		pool,pattern,sale_coef,floor_coef,total_coef,kitchen_coef,balcony_coef,metro_coef,state_coef) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
		$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24) returning id;`
	selectPatterns = `select pattern, avg_price from patterns where pool_id = $1;`
	selectAnalogs  = `select id,lng,lat,addr,room,segment,floors,cur_floor,walls,total,kitchen,balcony,metro,state,price,avg_price,
	sale_coef,floor_coef,total_coef,kitchen_coef,balcony_coef,metro_coef,state_coef from analogs where use = 't' and pool = $1 and pattern = $2;`
	insertCookie = `insert into cookies(user_id, cookie) values($1, $2);`
	deleteCookie = `delete from cookies where cookie = $1;`
	selectCookie = `select user_id from cookies where cookie = $1;`
	insertUser   = `insert into users(login, pass) values($1, $2) returning id;`
	selectUser   = `select id, login, pass from users where login = $1;`
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

func (repo *dbRepository) SaveAnalogs(id int, pattern *domain.Row, analogs []*domain.Row, coefs []*domain.CorrectCoefs) error {
	ctx := context.Background()
	_, err := repo.db.Pool.Exec(ctx, insertPattern, id, pattern.ID, pattern.Longitude, pattern.Latitude, pattern.AvgCost)
	if err != nil {
		log.Error("Unable to save pattern: ", err)
		return err
	}
	wg := sync.WaitGroup{}
	for i, a := range analogs {
		wg.Add(1)
		go func(i int, a *domain.Row) {
			defer wg.Done()
			err := repo.db.Pool.QueryRow(ctx, insertAnalog, a.Longitude, a.Latitude, a.Address, a.Rooms, a.Segment, a.Floors, a.CFloor, a.Walls, a.Total,
				a.Kitchen, a.Balcony, a.Metro, a.State, a.Cost, a.AvgCost, id, pattern.ID, coefs[i].Sale, coefs[i].Floor, coefs[i].Total, coefs[i].Kitchen,
				coefs[i].Balcony, coefs[i].Metro, coefs[i].State).Scan(&a.ID)
			if err != nil {
				log.Error("Unable to save analog: ", err)
			}
		}(i, a)
	}
	wg.Wait()
	return nil
}

func (repo *dbRepository) GetPatternAnalogs(id int) ([]*domain.PatternAnalogs, error) {
	rows, err := repo.db.Pool.Query(context.Background(), selectPatterns, id)
	if err != nil {
		log.Error(fmt.Sprint("Cannot get patterns:", id, err))
		return nil, err
	}
	defer rows.Close()

	result := make([]*domain.PatternAnalogs, 0)
	for rows.Next() {
		ptn := &domain.PatternAnalogs{Pattern: new(domain.Row)}
		err = rows.Scan(&ptn.Pattern.ID, &ptn.Pattern.AvgCost)
		if err != nil {
			continue
		}
		ptn.Analogs, ptn.Correct = repo.GetAnalogs(id, ptn.Pattern.ID)
		result = append(result, ptn)
	}
	return result, nil
}

func (repo *dbRepository) GetAnalogs(id, ptnID int) ([]*domain.Row, []*domain.CorrectCoefs) {
	rows, err := repo.db.Pool.Query(context.Background(), selectAnalogs, id, ptnID)
	if err != nil {
		log.Error(fmt.Sprint("Cannot get pattern analogs:", id, ptnID, err))
		return nil, nil
	}
	defer rows.Close()

	analogs, coefs := make([]*domain.Row, 0), make([]*domain.CorrectCoefs, 0)
	for rows.Next() {
		a, c := new(domain.Row), new(domain.CorrectCoefs)
		err = rows.Scan(&a.ID, &a.Longitude, &a.Latitude, &a.Address, &a.Rooms, &a.Segment, &a.Floors, &a.CFloor, &a.Walls, &a.Total, &a.Kitchen,
			&a.Balcony, &a.Metro, &a.State, &a.Cost, &a.AvgCost, &c.Sale, &c.Floor, &c.Total, &c.Kitchen, &c.Balcony, &c.Metro, &c.State)
		if err != nil {
			continue
		}
		analogs, coefs = append(analogs, a), append(coefs, c)
	}
	return analogs, coefs
}

func (repo *dbRepository) CreateCookie(id int, cookie string) error {
	_, err := repo.db.Pool.Exec(context.Background(), insertCookie, id, cookie)
	if err != nil {
		log.Info("Cookie creating: ", err)
	}
	return err
}

func (repo *dbRepository) DeleteCookie(cookie string) {
	_, err := repo.db.Pool.Exec(context.Background(), deleteCookie, cookie)
	if err != nil {
		log.Info("Cookie deleting: ", err)
	}
}

func (repo *dbRepository) CheckCookie(cookie string) int {
	user := 0
	err := repo.db.Pool.QueryRow(context.Background(), selectCookie, cookie).Scan(&user)
	if err != nil {
		log.Info("Cookie checking: ", err)
	}
	return user
}

func (repo *dbRepository) CreateUser(user *domain.User) error {
	err := repo.db.Pool.QueryRow(context.Background(), insertUser, user.Login, user.HashPass).Scan(&user.ID)
	if err != nil {
		log.Error("Cannot create user: ", err)
	}
	return err
}

func (repo *dbRepository) GetUser(login string) *domain.User {
	user := new(domain.User)
	err := repo.db.Pool.QueryRow(context.Background(), selectUser, login).Scan(&user.ID, &user.Login, &user.HashPass)
	if err != nil {
		log.Error("Cannot get user: ", err)
		return nil
	}
	return user
}

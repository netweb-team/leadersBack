package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/domain"
	"leaders_apartments/internal/pkg/parser/html"
	"leaders_apartments/internal/pkg/parser/xslx"
	"net/http"
	"os"

	"github.com/labstack/gommon/log"
)

type serverUsecases struct {
	repo domain.Repository
}

func New(repo domain.Repository) domain.Usecase {
	return &serverUsecases{repo}
}

func (u *serverUsecases) ImportXslx(f io.Reader) *domain.Table {
	data, err := xslx.Parse(f)
	if err != nil {
		return nil
	}

	if data.ID, err = u.repo.SaveTable(data.Path); err != nil {
		os.Remove(data.Path)
		return nil
	}
	return data
}

func (u *serverUsecases) FindAnalogs(id, ptnIndex int) *domain.PatternAnalogs {
	filename, err := u.repo.GetTableName(id)
	if err != nil {
		return nil
	}
	row := xslx.ReadRow(filename, ptnIndex)
	if row == nil {
		return nil
	}
	row.ID = ptnIndex
	resp, err := http.Get(fmt.Sprintf(config.New().MapApi, row.Address, config.New().ApiKey))
	if err != nil {
		log.Info("Cannot get request for coordinates: ", err)
		return nil
	}
	defer resp.Body.Close()
	coords := new(domain.MapResponse)
	json.NewDecoder(resp.Body).Decode(coords)
	row.Longitude, row.Latitude = coords.Features[0].Center[0], coords.Features[0].Center[1]
	analogs := html.FindAnalogs(row)
	if len(analogs) == 0 {
		return nil
	}
	for _, anal := range analogs {
		anal.Rooms, anal.Floors, anal.Walls = row.Rooms, row.Floors, row.Walls
	}
	// save something in db
	return &domain.PatternAnalogs{Pattern: row, Analogs: analogs}
}

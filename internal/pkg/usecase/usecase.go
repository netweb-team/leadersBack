package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/correct"
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

func (u *serverUsecases) ImportXlsx(f io.Reader) *domain.Table {
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

func (u *serverUsecases) ExportXlsx(id int) string {
	filename, err := u.repo.GetTableName(id)
	if err != nil {
		return ""
	}
	return filename
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
	cfg := config.New()
	resp, err := http.Get(fmt.Sprintf(cfg.MapApi, row.Address, cfg.ApiKey))
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
		anal.Rooms, anal.Floors, anal.Walls, anal.Good, anal.AvgCost = row.Rooms, row.Floors, row.Walls, 1, float64(anal.Cost)/anal.Total
	}
	coefs := correct.Do(row, analogs)
	sum := 0.0
	for _, anal := range analogs {
		sum += anal.AvgCost
	}
	row.AvgCost = sum / float64(len(analogs))
	row.Cost = int(row.Total * row.AvgCost)
	if u.repo.SaveAnalogs(id, row, analogs, coefs) != nil {
		return nil
	}
	return &domain.PatternAnalogs{Pattern: row, Analogs: analogs, Correct: coefs}
}

func (u *serverUsecases) CalcPool(id int) []*domain.Row {
	filename, err := u.repo.GetTableName(id)
	if err != nil {
		return nil
	}
	table := xslx.ReadTable(filename)
	if table == nil {
		return nil
	}
	data, err := u.repo.GetPatternAnalogs(id)
	if err != nil {
		return nil
	}
	for i := range data {
		idx := data[i].Pattern.ID - 1
		table[idx].AvgCost = data[i].Pattern.AvgCost
		data[i].Pattern = table[idx]
	}
	groups := make([][]*domain.Row, len(data))
	for i := range groups {
		groups[i] = make([]*domain.Row, 0)
	}
	for j, row := range table {
		for i, ptn := range data {
			p := ptn.Pattern
			if j == p.ID-1 || row.Rooms == p.Rooms && row.Segment == p.Segment && row.Floors == p.Floors && row.Walls == p.Walls {
				row.AvgCost = p.AvgCost
				groups[i] = append(groups[i], row)
				break
			}
		}
	}
	for i := range groups {
		correct.Do(data[i].Pattern, groups[i])
	}
	_ = xslx.SavePrice(filename, table)
	//save to db archive
	return table
}

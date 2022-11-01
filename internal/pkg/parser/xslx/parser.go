package xslx

import (
	"errors"
	"io"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/domain"
	"leaders_apartments/internal/pkg/utils"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/plandem/xlsx"
)

const (
	ext      = ".xslx"
	nameLen  = 16
	segments = "новостройка;современное жилье;старый жилой фонд"
	walls    = "кирпич;панель;монолит"
	states   = "без отделки;муниципальный ремонт;современная отделка"
)

func Parse(f io.Reader) (*domain.Table, error) {
	xl, err := xlsx.Open(f)
	if err != nil {
		log.Error("Cannot open file: ", err)
		return nil, err
	}
	defer xl.Close()
	if err := xl.IsValid(); err != nil {
		log.Error("Xslx is invalid: ", err)
		return nil, err
	}

	sheet := xl.Sheet(0)
	totalCols, totalRows := sheet.Dimension()
	result := &domain.Table{Rows: make([]*domain.Row, 0)}
	for rIdx := 1; rIdx < totalRows; rIdx++ {
		row := new(domain.Row)
		for cIdx := 0; cIdx < totalCols; cIdx++ {
			c := sheet.Cell(rIdx, cIdx)
			switch cIdx {
			case 0:
				row.Address = c.Value()
			case 1:
				if row.Rooms, err = c.Uint(); err != nil {
					log.Info("Room count in table is not uint: ", err)
					return nil, err
				}
			case 2:
				if !strings.Contains(segments, strings.ToLower(c.Value())) {
					log.Info("Unknown segment of building")
					return nil, errors.New("Unknown segment of building")
				}
				row.Segment = c.Value()
			case 3:
				if row.Floors, err = c.Uint(); err != nil {
					log.Info("Total floors in table is not uint: ", err)
					return nil, err
				}
			case 4:
				if !strings.Contains(walls, strings.ToLower(c.Value())) {
					log.Info("Unknown material of walls")
					return nil, errors.New("Unknown material of walls")
				}
				row.Walls = c.Value()
			case 5:
				if row.CFloor, err = c.Uint(); err != nil {
					log.Info("Current floor in table is not uint: ", err)
					return nil, err
				}
			case 6:
				if row.Total, err = c.Float(); err != nil {
					log.Info("Total square in table is not float: ", err)
					return nil, err
				}
			case 7:
				if row.Kitchen, err = c.Float(); err != nil {
					log.Info("Kitchen square in table is not float: ", err)
					return nil, err
				}
			case 8:
				row.Balcony = c.Value()
			case 9:
				row.Metro = c.Value()
			case 10:
				if !strings.Contains(states, strings.ToLower(c.Value())) {
					log.Info("Unknown state of flat")
					return nil, errors.New("Unknown state of flat")
				}
				row.State = c.Value()
			}
		}
		result.Rows = append(result.Rows, row)
	}
	result.Path = config.New().Path + utils.RandString(nameLen) + ext
	if err := xl.SaveAs(result.Path); err != nil {
		log.Error("Cannot save file to path:", result.Path, err)
		return nil, err
	}
	return result, nil
}

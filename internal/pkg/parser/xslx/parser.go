package xslx

import (
	"errors"
	"fmt"
	"io"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/domain"
	"leaders_apartments/internal/pkg/utils"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/plandem/xlsx"
)

const (
	ext     = ".xlsx"
	nameLen = 16
)

var (
	segments = fmt.Sprintf("%s;%s;%s", domain.LowerSegmentNew, domain.LowerSegmentMid, domain.LowerSegmentOld)
	walls    = fmt.Sprintf("%s;%s;%s", domain.LowerWallBrick, domain.LowerWallMono, domain.LowerWallPanel)
	states   = fmt.Sprintf("%s;%s;%s", domain.LowerStateOff, domain.LowerStateMun, domain.LowerStateNew)
)

func Parse(f io.Reader) (*domain.Table, error) {
	xl, err := xlsx.Open(f)
	if err != nil {
		log.Error("Cannot open file for import: ", err)
		return nil, err
	}
	defer xl.Close()
	if err := xl.IsValid(); err != nil {
		log.Error("Xslx is invalid: ", err)
		return nil, err
	}

	sheet := xl.Sheet(0)
	result := new(domain.Table)
	if table := parseSheet(sheet); len(table) == 0 {
		return nil, errors.New("Error while parsing xlsx table")
	} else {
		result.Rows = table
		result.Flat = len(table)
	}
	result.Path = config.New().Path + utils.RandString(nameLen) + ext
	if err := xl.SaveAs(result.Path); err != nil {
		log.Error("Cannot save file to path:", result.Path, err)
		return nil, err
	}
	return result, nil
}

func ReadTable(filename string) []*domain.Row {
	xl, err := xlsx.Open(filename)
	if err != nil {
		log.Error("Cannot open file for read: ", err)
		return nil
	}
	defer xl.Close()
	return parseSheet(xl.Sheet(0))
}

func ReadRow(filename string, id int) *domain.Row {
	xl, err := xlsx.Open(filename)
	if err != nil {
		log.Error("Cannot open file for read: ", err)
		return nil
	}
	defer xl.Close()

	sheet := xl.Sheet(0)
	cols, rows := sheet.Dimension()
	if id <= 0 || id >= rows {
		log.Info("Incorrect number of row: ", id)
		return nil
	}
	return parseRow(sheet, id, cols)
}

func parseSheet(sheet xlsx.Sheet) []*domain.Row {
	result := make([]*domain.Row, 0)
	totalCols, totalRows := sheet.Dimension()
	for rIdx := 1; rIdx < totalRows; rIdx++ {
		if row := parseRow(sheet, rIdx, totalCols); row != nil {
			result = append(result, row)
		} else {
			break
		}
	}
	return result
}

func parseRow(sheet xlsx.Sheet, i, total int) *domain.Row {
	var err error
	row := new(domain.Row)
	for col := 0; col < total; col++ {
		c := sheet.Cell(col, i)
		switch col {
		case 0:
			row.Address = c.Value()
		case 1:
			if _, err = c.Uint(); err != nil && strings.ToLower(c.Value()) != domain.LowerStudio {
				log.Info("Room count in table is not uint or studio: ", err)
				return nil
			}
			row.Rooms = c.Value()
		case 2:
			if !strings.Contains(segments, strings.ToLower(c.Value())) {
				log.Info("Unknown segment of building")
				return nil
			}
			row.Segment = c.Value()
		case 3:
			if row.Floors, err = c.Uint(); err != nil {
				log.Info("Total floors in table is not uint: ", err)
				return nil
			}
		case 4:
			if !strings.Contains(walls, strings.ToLower(c.Value())) {
				log.Info("Unknown material of walls")
				return nil
			}
			row.Walls = c.Value()
		case 5:
			if row.CFloor, err = c.Uint(); err != nil {
				log.Info("Current floor in table is not uint: ", err)
				return nil
			}
		case 6:
			if row.Total, err = c.Float(); err != nil {
				log.Info("Total square in table is not float: ", err)
				return nil
			}
		case 7:
			if row.Kitchen, err = c.Float(); err != nil {
				log.Info("Kitchen square in table is not float: ", err)
				return nil
			}
		case 8:
			b := strings.ToLower(c.Value())
			if b != strings.ToLower(domain.Yes) && b != domain.LowerNo {
				log.Info("Balcony is not yes/no")
				return nil
			}
			row.Balcony = c.Value()
		case 9:
			if row.Metro, err = c.Float(); err != nil {
				log.Info("Metro in table is not float: ", err)
				return nil
			}
		case 10:
			if !strings.Contains(states, strings.ToLower(c.Value())) {
				log.Info("Unknown state of flat")
				return nil
			}
			row.State = c.Value()
		}
	}
	return row
}

func SavePrice(filename string, table []*domain.Row) string {
	xl, err := xlsx.Open(filename)
	if err != nil {
		log.Error("Cannot open file for read: ", err)
		return ""
	}
	defer xl.Close()
	sh := xl.Sheet(0)

	xlp := xlsx.New()
	defer xlp.Close()
	shp := xlp.AddSheet("Sheet 1")
	cols, rows := sh.Dimension()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			shp.Cell(j, i).SetValue(sh.Cell(j, i))
		}
	}
	shp.CellByRef("L1").SetText("Цена")
	for i, row := range table {
		row.Cost = int(row.AvgCost * row.Total)
		shp.Cell(11, i+1).SetValue(row.Cost)
	}

	path := config.New().Path + utils.RandString(nameLen) + ext
	xlp.SaveAs(path)
	return path
}

func ReadPrice(filename string, table []*domain.Row) error {
	xl, err := xlsx.Open(filename)
	if err != nil {
		log.Error("Cannot open file for read: ", err)
		return err
	}
	defer xl.Close()
	sh := xl.Sheet(0)

	for i, row := range table {
		c := sh.Cell(11, i+1)
		if row.Cost, err = c.Int(); err != nil {
			log.Error("Cannot read price: ", err)
		}
	}
	return nil
}

package usecase

import (
	"io"
	"leaders_apartments/internal/pkg/domain"
	"leaders_apartments/internal/pkg/parser/xslx"
	"os"
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

func (u *serverUsecases) FindAnalogs(id, ptnIndex int) []*domain.AdPage {
	// get pattern from xslx
	// make mapbox req for coordinates
	// find analogs cian
	// save something in db
	return nil
}

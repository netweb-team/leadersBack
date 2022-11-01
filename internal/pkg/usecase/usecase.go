package usecase

import (
	"leaders_apartments/internal/pkg/domain"
	"os"
)

type serverUsecases struct {
	repo domain.Repository
}

func New(repo domain.Repository) domain.Usecase {
	return &serverUsecases{repo}
}

func (u *serverUsecases) ImportXslx(f *os.File) error {
	return nil
}

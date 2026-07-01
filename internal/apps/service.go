package apps

import (
	"mobile-connect/internal/adb"
	"mobile-connect/internal/apps/entity"
)

type Service struct {
	adbClient  adb.ADBClient
	repository *Repository
}

func NewService(adbClient adb.ADBClient, repository *Repository) *Service {
	return &Service{
		adbClient:  adbClient,
		repository: repository,
	}
}

func (s *Service) Create(app entity.App) error {
	return s.repository.Create(app)
}

func (s *Service) FindAll() ([]entity.App, error) {
	return s.repository.FindAll()
}

func (s *Service) FindByID(id string) (*entity.App, error) {
	return s.repository.FindByID(id)
}

func (s *Service) Delete(id string) error {
	return s.repository.Delete(id)
}

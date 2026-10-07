package apps

import (
	"mobile-connect/internal/adb"
	"mobile-connect/internal/apk"
	"mobile-connect/internal/apps/entity"
	"os"
	"path/filepath"
)

type Service struct {
	adbClient  adb.ADBClient
	repository *Repository
	apkService *apk.Service
}

func NewService(adbClient adb.ADBClient, repository *Repository, apkService *apk.Service) *Service {
	return &Service{
		adbClient:  adbClient,
		repository: repository,
		apkService: apkService,
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

func resolveAPKPath(fileName string) string {
	storagePath := filepath.Join("storage", "apks", fileName)
	if _, err := os.Stat(storagePath); err == nil {
		return storagePath
	}

	return filepath.Join("uploads", fileName)
}

func (s *Service) Install(
	appID string,
	deviceID string,
) (string, error) {

	app, err := s.repository.FindByID(appID)
	if err != nil {
		return "", err
	}

	apkPath := resolveAPKPath(app.FileName)

	return s.apkService.InstallAPK(
		deviceID,
		apkPath,
	)
}

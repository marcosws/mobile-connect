package apk

import (
	"mobile-connect/internal/adb"
	"mobile-connect/internal/logger"
)

type Service struct {
	adb adb.ADBClient
}

func NewService(client adb.ADBClient) *Service {
	return &Service{
		adb: client,
	}
}

func (s *Service) InstallAPK(
	id string,
	apkPath string,
) (string, error) {

	logger.Logger.Printf(
		"installing apk device=%s file=%s",
		id,
		apkPath,
	)

	output, err := s.adb.InstallAPK(
		id,
		apkPath,
	)

	if err != nil {

		logger.Logger.Printf(
			"install apk failed device=%s error=%v",
			id,
			err,
		)

		return output, err
	}

	logger.Logger.Printf(
		"install apk success device=%s",
		id,
	)

	return output, nil
}

func (s *Service) UninstallAPK(
	id string,
	packageName string,
) (string, error) {

	logger.Logger.Printf(
		"uninstalling package device=%s package=%s",
		id,
		packageName,
	)

	output, err := s.adb.UninstallAPK(
		id,
		packageName,
	)

	if err != nil {

		logger.Logger.Printf(
			"uninstall failed device=%s package=%s error=%v",
			id,
			packageName,
			err,
		)

		return output, err
	}

	logger.Logger.Printf(
		"uninstall success device=%s package=%s",
		id,
		packageName,
	)

	return output, nil
}
func (s *Service) ListPackages(
	id string,
) ([]string, error) {

	logger.Logger.Printf(
		"listing packages device=%s",
		id,
	)

	packages, err := s.adb.ListPackages(id)

	if err != nil {

		logger.Logger.Printf(
			"list packages failed device=%s error=%v",
			id,
			err,
		)

		return nil, err
	}

	logger.Logger.Printf(
		"packages found device=%s total=%d",
		id,
		len(packages),
	)

	return packages, nil
}

func (s *Service) LaunchApp(
	deviceID string,
	pkg string,
) error {

	return s.adb.LaunchApp(
		deviceID,
		pkg,
	)
}

func (s *Service) StopApp(
	deviceID string,
	pkg string,
) error {

	return s.adb.StopApp(
		deviceID,
		pkg,
	)
}

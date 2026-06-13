package shell

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

func (s *Service) Shell(
	id string,
	command string,
) (string, error) {

	logger.Logger.Printf(
		"executing shell device=%s command=%s",
		id,
		command,
	)

	output, err := s.adb.Shell(
		id,
		command,
	)

	if err != nil {

		logger.Logger.Printf(
			"shell failed device=%s error=%v",
			id,
			err,
		)

		return "", err
	}

	logger.Logger.Printf(
		"shell success device=%s",
		id,
	)

	return output, nil
}

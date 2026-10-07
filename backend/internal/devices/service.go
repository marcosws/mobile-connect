package devices

import (
	"mobile-connect/internal/adb"
	"mobile-connect/internal/devices/entity"
	"mobile-connect/internal/logger"
	"strings"
)

type Service struct {
	adb adb.ADBClient
}

func NewService(client adb.ADBClient) *Service {
	return &Service{
		adb: client,
	}
}

func (s *Service) List() ([]entity.Device, error) {

	logger.Logger.Printf(
		"listing devices",
	)

	output, err := s.adb.Devices()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")

	var devices []entity.Device

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" ||
			strings.HasPrefix(line, "List of devices") {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		serial := fields[0]

		manufacturer, _ := s.adb.GetProp(
			serial,
			"ro.product.manufacturer",
		)

		model, _ := s.adb.GetProp(
			serial,
			"ro.product.model",
		)

		androidVersion, _ := s.adb.GetProp(
			serial,
			"ro.build.version.release",
		)

		sdk, _ := s.adb.GetProp(
			serial,
			"ro.build.version.sdk",
		)

		devices = append(devices, entity.Device{
			ID:             serial,
			Status:         fields[1],
			Manufacturer:   manufacturer,
			Model:          model,
			AndroidVersion: androidVersion,
			SDK:            sdk,
		})
		logger.Logger.Printf(
			"device found: %s manufacturer=%s model=%s androidVersion=%s sdk=%s",
			devices[len(devices)-1].ID,
			devices[len(devices)-1].Manufacturer,
			devices[len(devices)-1].Model,
			devices[len(devices)-1].AndroidVersion,
			devices[len(devices)-1].SDK,
		)
	}

	return devices, nil
}

func (s *Service) GetByID(id string) (*entity.Device, error) {

	devices, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.ID == id {
			return &device, nil
		}
	}

	return nil, nil
}

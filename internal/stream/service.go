package stream

import (
	"os/exec"
	"strconv"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Capture(deviceID string) ([]byte, error) {

	cmd := exec.Command(
		"adb",
		"-s",
		deviceID,
		"exec-out",
		"screencap",
		"-p",
	)

	return cmd.Output()
}

func (s *Service) Tap(
	deviceID string,
	x int,
	y int,
) error {

	cmd := exec.Command(
		"adb",
		"-s",
		deviceID,
		"shell",
		"input",
		"tap",
		strconv.Itoa(x),
		strconv.Itoa(y),
	)

	return cmd.Run()
}

func (s *Service) Swipe(
	deviceID string,
	x1 int,
	y1 int,
	x2 int,
	y2 int,
	duration int,
) error {

	cmd := exec.Command(
		"adb",
		"-s",
		deviceID,
		"shell",
		"input",
		"swipe",
		strconv.Itoa(x1),
		strconv.Itoa(y1),
		strconv.Itoa(x2),
		strconv.Itoa(y2),
		strconv.Itoa(duration),
	)

	return cmd.Run()
}

func (s *Service) LongPress(
	deviceID string,
	x int,
	y int,
	duration int,
) error {

	cmd := exec.Command(
		"adb",
		"-s",
		deviceID,
		"shell",
		"input",
		"swipe",
		strconv.Itoa(x),
		strconv.Itoa(y),
		strconv.Itoa(x),
		strconv.Itoa(y),
		strconv.Itoa(duration),
	)

	return cmd.Run()
}

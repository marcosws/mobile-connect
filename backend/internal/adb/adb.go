package adb

import (
	"bytes"
	"os/exec"
	"strings"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Devices() (string, error) {
	cmd := exec.Command("adb", "devices")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (c *Client) GetProp(serial, prop string) (string, error) {

	cmd := exec.Command(
		"adb",
		"-s",
		serial,
		"shell",
		"getprop",
		prop,
	)

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func (c *Client) Shell(serial string, command string) (string, error) {

	args := []string{
		"-s",
		serial,
		"shell",
	}

	args = append(args, strings.Fields(command)...)

	cmd := exec.Command(
		"adb",
		args...,
	)

	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (c *Client) InstallAPK(
	serial string,
	apkPath string,
) (string, error) {

	cmd := exec.Command(
		"adb",
		"-s",
		serial,
		"install",
		"-r",
		apkPath,
	)

	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (c *Client) UninstallAPK(
	serial string,
	packageName string,
) (string, error) {

	cmd := exec.Command(
		"adb",
		"-s",
		serial,
		"uninstall",
		packageName,
	)

	output, err := cmd.CombinedOutput()

	return string(output), err
}
func (c *Client) ListPackages(
	serial string,
) ([]string, error) {

	cmd := exec.Command(
		"adb",
		"-s",
		serial,
		"shell",
		"pm",
		"list",
		"packages",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(
		string(output),
		"\n",
	)

	var packages []string

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		line = strings.TrimPrefix(
			line,
			"package:",
		)

		packages = append(
			packages,
			line,
		)
	}

	return packages, nil
}

func (c *Client) LaunchApp(
	deviceID string,
	pkg string,
) error {

	cmd := exec.Command(
		"adb",
		"-s",
		deviceID,
		"shell",
		"monkey",
		"-p",
		pkg,
		"-c",
		"android.intent.category.LAUNCHER",
		"1",
	)

	return cmd.Run()
}

func (c *Client) StopApp(
	deviceID string,
	pkg string,
) error {

	cmd := exec.Command(
		"adb",
		"-s",
		deviceID,
		"shell",
		"am",
		"force-stop",
		pkg,
	)

	return cmd.Run()
}

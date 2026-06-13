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

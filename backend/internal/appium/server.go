package appium

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const startupTimeout = 20 * time.Second

// StartLocal starts Appium only when the configured endpoint is not already
// available. The process inherits stdout and stderr so Appium logs remain visible
// beside Mobile Connect logs. It is stopped when ctx is cancelled.
func StartLocal(ctx context.Context, client *Client, command string) error {
	checkCtx, cancel := context.WithTimeout(ctx, time.Second)
	alreadyRunning := client.Available(checkCtx)
	cancel()
	if alreadyRunning {
		return nil
	}

	endpoint := client.BaseURL()
	host := endpoint.Hostname()
	if !isLocalHost(host) {
		return fmt.Errorf("cannot auto-start Appium for non-local APPIUM_URL host %q", host)
	}
	port := endpoint.Port()
	if port == "" {
		return fmt.Errorf("APPIUM_URL must include an explicit port when auto-starting Appium")
	}

	args := []string{"--address", host, "--port", port}
	if endpoint.Path != "" && endpoint.Path != "/" {
		args = append(args, "--base-path", endpoint.Path)
	}
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start Appium command %q: %w", command, err)
	}

	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()

	deadline := time.NewTimer(startupTimeout)
	defer deadline.Stop()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-exited:
			return fmt.Errorf("Appium exited during startup: %w", err)
		case <-deadline.C:
			_ = cmd.Process.Kill()
			<-exited
			return fmt.Errorf("Appium did not become ready within %s", startupTimeout)
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			checkCtx, cancel := context.WithTimeout(ctx, time.Second)
			available := client.Available(checkCtx)
			cancel()
			if available {
				return nil
			}
		}
	}
}

func isLocalHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

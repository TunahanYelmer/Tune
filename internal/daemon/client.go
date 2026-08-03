package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func SocketPath() string {
	return filepath.Join(os.TempDir(), "tune-daemon.sock")
}

func LogPath() string {
	return filepath.Join(os.TempDir(), "tune-daemon.log")
}

// EnsureRunning connects to the daemon, starting it in the background
// if it isn't already running.
func EnsureRunning() error {
	if tryConnect() == nil {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	logFile, err := os.Create(LogPath())
	if err != nil {
		return fmt.Errorf("creating daemon log file: %w", err)
	}

	cmd := exec.Command(exe, "__daemon")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting daemon: %w", err)
	}

	for i := 0; i < 30; i++ {
		if tryConnect() == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("daemon did not start in time (check %s)", LogPath())
}

func tryConnect() error {
	conn, err := net.DialTimeout("unix", SocketPath(), 200*time.Millisecond)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// Send issues a request to the running daemon and returns its response.
func Send(req Request) (*Response, error) {
	conn, err := net.Dial("unix", SocketPath())
	if err != nil {
		return nil, fmt.Errorf("connecting to daemon: %w", err)
	}
	defer conn.Close()

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(append(data, '\n')); err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	reader := bufio.NewReader(conn)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}
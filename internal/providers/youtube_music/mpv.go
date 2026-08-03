package youtube_music

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type mpvClient struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	sockPath string
}

func socketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\tune-mpv`
	}

	return filepath.Join(os.TempDir(), "tune-mpv.sock")
}

func mpvLogPath() string {
	return filepath.Join(os.TempDir(), "tune-mpv.log")
}

// newMPVClient launches mpv in idle mode with an IPC socket.
func newMPVClient() (*mpvClient, error) {
	sock := socketPath()

// Reuse existing mpv if it is already running
if _, err := os.Stat(sock); err == nil {

	conn, err := net.Dial("unix", sock)

	if err == nil {
		conn.Close()

		return &mpvClient{
			sockPath: sock,
		}, nil
	}

	// Socket exists but mpv is dead
	os.Remove(sock)
}

	logFile, err := os.Create(mpvLogPath())
	if err != nil {
		return nil, fmt.Errorf("creating mpv log file: %w", err)
	}

	cmd := exec.Command(
		"mpv",
		"--idle=yes",
		"--no-video",
		
		fmt.Sprintf("--input-ipc-server=%s", sock),
	)

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting mpv: %w (is mpv installed?)", err)
	}

	// Wait for mpv socket
	var errConn error

	for i := 0; i < 30; i++ {

		conn, err := net.Dial("unix", sock)

		if err == nil {
			conn.Close()
			errConn = nil
			break
		}

		errConn = err
		time.Sleep(100 * time.Millisecond)
	}

	if errConn != nil {
		return nil, fmt.Errorf(
			"connecting to mpv socket: %w (check %s)",
			errConn,
			mpvLogPath(),
		)
	}

	return &mpvClient{
		cmd:      cmd,
		sockPath: sock,
	}, nil
}


// send sends a command to mpv and waits for the response.
func (m *mpvClient) send(command []interface{}) (map[string]interface{}, error) {

	m.mu.Lock()
	defer m.mu.Unlock()


	conn, err := net.Dial("unix", m.sockPath)

	if err != nil {
		return nil, fmt.Errorf("connecting to mpv: %w", err)
	}

	defer conn.Close()


	payload, err := json.Marshal(map[string]interface{}{
		"command": command,
	})

	if err != nil {
		return nil, err
	}


	if _, err := conn.Write(append(payload, '\n')); err != nil {
		return nil, fmt.Errorf("writing to mpv: %w", err)
	}


	reader := bufio.NewReader(conn)


	for {

		line, err := reader.ReadBytes('\n')

		if err != nil {
			return nil, fmt.Errorf("reading mpv response: %w", err)
		}


		var resp map[string]interface{}

		if err := json.Unmarshal(line, &resp); err != nil {
			continue
		}


		// Ignore async events
		if _, ok := resp["event"]; ok {
			continue
		}


		return resp, nil
	}
}


// loadFile loads a stream and starts playback.
func (m *mpvClient) loadFile(url string) error {

	resp, err := m.send([]interface{}{
		"loadfile",
		url,
		"replace",
	})


	if err != nil {
		return err
	}


	if errStr, ok := resp["error"].(string); ok && errStr != "success" {
		return fmt.Errorf(
			"mpv loadfile error: %s",
			errStr,
		)
	}


	// Give mpv time to initialize playback
	time.Sleep(200 * time.Millisecond)


	// Force resume playback
	if err := m.setProperty("pause", false); err != nil {
		return fmt.Errorf(
			"unpausing mpv: %w",
			err,
		)
	}


	return nil
}


// setProperty changes an mpv property.
func (m *mpvClient) setProperty(
	name string,
	value interface{},
) error {

	resp, err := m.send([]interface{}{
		"set_property",
		name,
		value,
	})


	if err != nil {
		return err
	}


	if errStr, ok := resp["error"].(string); ok && errStr != "success" {
		return fmt.Errorf(
			"mpv set property error: %s",
			errStr,
		)
	}


	return nil
}


// getProperty gets an mpv property.
func (m *mpvClient) getProperty(
	name string,
) (interface{}, error) {

	resp, err := m.send([]interface{}{
		"get_property",
		name,
	})


	if err != nil {
		return nil, err
	}


	return resp["data"], nil
}


// togglePause toggles playback pause state.
func (m *mpvClient) togglePause() error {

	current, err := m.getProperty("pause")

	if err != nil {
		return err
	}


	isPaused, ok := current.(bool)

	if !ok {
		return fmt.Errorf("invalid pause state")
	}


	return m.setProperty(
		"pause",
		!isPaused,
	)
}
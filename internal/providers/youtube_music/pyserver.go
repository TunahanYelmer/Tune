package youtube_music

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/tunahanyelmer/Tune/internal/providers"
)

type pyServer struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
}

func newPyServer() (*pyServer, error) {
	cmd := exec.Command("python3", "internal/providers/youtube_music/scripts/server.py")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("opening python stdin: %w", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("opening python stdout: %w", err)
	}

	var stderrBuf []byte
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("opening python stderr: %w", err)
	}
	go func() {
		b, _ := io.ReadAll(stderrPipe)
		stderrBuf = b
	}()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting python search server: %w", err)
	}

	reader := bufio.NewReader(stdoutPipe)

	// Wait for the {"ready": true} line before returning, so callers
	// never race the ytmusicapi client's own init time.
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("python server failed to start: %w (stderr: %s)", err, string(stderrBuf))
	}
	var ready struct {
		Ready bool `json:"ready"`
	}
	if err := json.Unmarshal(line, &ready); err != nil || !ready.Ready {
		return nil, fmt.Errorf("unexpected python server startup output: %s", string(line))
	}

	return &pyServer{cmd: cmd, stdin: stdin, stdout: reader}, nil
}

func (s *pyServer) search(query string, limit int) ([]internalTrack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, err := json.Marshal(map[string]interface{}{"query": query, "limit": limit})
	if err != nil {
		return nil, err
	}
	if _, err := s.stdin.Write(append(req, '\n')); err != nil {
		return nil, fmt.Errorf("writing to python server: %w", err)
	}

	line, err := s.stdout.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("reading from python server: %w", err)
	}

	var resp struct {
		Tracks []struct {
			Title    string `json:"title"`
			Artist   string `json:"artist"`
			Album    string `json:"album"`
			Duration string `json:"duration"`
			VideoID  string `json:"videoId"`
		} `json:"tracks"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("parsing python response: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("ytmusicapi error: %s", resp.Error)
	}

	tracks := make([]internalTrack, len(resp.Tracks))
	for i, r := range resp.Tracks {
		tracks[i] = internalTrack{
			Track: providers.Track{
				Title:    r.Title,
				Artist:   r.Artist,
				Album:    r.Album,
				Duration: parseDuration(r.Duration),
			},
			VideoID: r.VideoID,
		}
	}
	return tracks, nil
}
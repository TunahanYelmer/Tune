package youtube_music

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/tunahanyelmer/Tune/internal/providers"
)

type internalTrack struct {
	providers.Track
	VideoID string
}

type YouTubeMusicProvider struct {
	mu      sync.Mutex
	mpv     *mpvClient
	py      *pyServer
	queue   []internalTrack
	index   int
	current *internalTrack
}

func NewYouTubeMusicProvider() *YouTubeMusicProvider {
	return &YouTubeMusicProvider{}
}

// --- Provider interface ---

func (p *YouTubeMusicProvider) Login() error {
	for _, bin := range []string{"python3", "yt-dlp", "mpv"} {
		if _, err := exec.LookPath(bin); err != nil {
			return fmt.Errorf(
				"%s not found in PATH — required for youtube_music provider",
				bin,
			)
		}
	}

	if err := exec.Command(
		"python3",
		"-c",
		"import ytmusicapi",
	).Run(); err != nil {
		return fmt.Errorf(
			"ytmusicapi not installed — run: pip install ytmusicapi",
		)
	}

	return nil
}


func (p *YouTubeMusicProvider) Search(query string) ([]providers.Track, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.ensurePy(); err != nil {
		return nil, err
	}

	internal, err := p.py.search(query, 5)
	if err != nil {
		return nil, err
	}

	tracks := make([]providers.Track, len(internal))

	for i, t := range internal {
		tracks[i] = t.Track
	}

	return tracks, nil
}


func (p *YouTubeMusicProvider) Play(query string) error {

	fmt.Println("PROVIDER PLAY:", query)

	p.mu.Lock()
	defer p.mu.Unlock()


	if err := p.ensureMPV(); err != nil {
		return err
	}

	if err := p.ensurePy(); err != nil {
		return err
	}


	results, err := p.py.search(query, 5)

	if err != nil {
		return err
	}


	if len(results) == 0 {
		return fmt.Errorf(
			"no results found for %q",
			query,
		)
	}


	p.queue = results
	p.index = 0


	return p.playCurrentLocked()
}


func (p *YouTubeMusicProvider) Pause() error {

	p.mu.Lock()
	defer p.mu.Unlock()


	if p.mpv == nil {
		return fmt.Errorf(
			"nothing is playing",
		)
	}


	return p.mpv.togglePause()
}



func (p *YouTubeMusicProvider) Next() error {

	p.mu.Lock()
	defer p.mu.Unlock()


	if p.index+1 >= len(p.queue) {
		return fmt.Errorf(
			"no next track queued",
		)
	}


	p.index++

	return p.playCurrentLocked()
}



func (p *YouTubeMusicProvider) Previous() error {

	p.mu.Lock()
	defer p.mu.Unlock()


	if p.index-1 < 0 {
		return fmt.Errorf(
			"already at first track",
		)
	}


	p.index--

	return p.playCurrentLocked()
}



func (p *YouTubeMusicProvider) Current() (*providers.Track, error) {

	p.mu.Lock()
	defer p.mu.Unlock()


	if p.current == nil {
		return nil, fmt.Errorf(
			"nothing is playing",
		)
	}


	t := p.current.Track

	return &t, nil
}



func (p *YouTubeMusicProvider) SetVolume(level int) error {

	p.mu.Lock()
	defer p.mu.Unlock()


	if p.mpv == nil {
		return fmt.Errorf(
			"nothing is playing",
		)
	}


	if level < 0 || level > 100 {
		return fmt.Errorf(
			"volume must be between 0 and 100",
		)
	}


	return p.mpv.setProperty(
		"volume",
		level,
	)
}



// --- internals ---


func (p *YouTubeMusicProvider) ensureMPV() error {

	if p.mpv != nil {
		return nil
	}


	client, err := newMPVClient()

	if err != nil {
		return err
	}


	p.mpv = client

	return nil
}



func (p *YouTubeMusicProvider) ensurePy() error {

	if p.py != nil {
		return nil
	}


	server, err := newPyServer()

	if err != nil {
		return err
	}


	p.py = server

	return nil
}



// playCurrentLocked assumes p.mu is already held.
func (p *YouTubeMusicProvider) playCurrentLocked() error {

	track := p.queue[p.index]


	streamURL, err := resolveStreamURL(
		track.VideoID,
	)

	if err != nil {
		return err
	}


	if err := p.mpv.loadFile(streamURL); err != nil {
		return fmt.Errorf(
			"playing track: %w",
			err,
		)
	}


	// loadFile already resumes playback


	p.current = &track


	return nil
}



func resolveStreamURL(videoID string) (string, error) {

	url := fmt.Sprintf(
		"https://music.youtube.com/watch?v=%s",
		videoID,
	)


	cmd := exec.Command(
		"yt-dlp",
		"-f",
		"bestaudio",
		"-g",
		url,
	)


	var stdout, stderr bytes.Buffer


	cmd.Stdout = &stdout
	cmd.Stderr = &stderr


	if err := cmd.Run(); err != nil {

		return "",
			fmt.Errorf(
				"resolving stream: %w (%s)",
				err,
				stderr.String(),
			)
	}


	return strings.TrimSpace(stdout.String()), nil
}



// parseDuration converts "3:08" -> 188 seconds.
func parseDuration(s string) int {

	parts := strings.Split(
		s,
		":",
	)


	if len(parts) != 2 {
		return 0
	}


	min, err1 := strconv.Atoi(parts[0])
	sec, err2 := strconv.Atoi(parts[1])


	if err1 != nil || err2 != nil {
		return 0
	}


	return min*60 + sec
}
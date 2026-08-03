package spotify

import (
	"fmt"

	"github.com/tunahanyelmer/Tune/internal/providers"
)

type Provider struct{}

func NewSpotifyProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Login() error {
	return fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) Play(query string) error {
	return fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) Pause() error {
	return fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) Next() error {
	return fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) Previous() error {
	return fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) Current() (*providers.Track, error) {
	return nil, fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) Search(query string) ([]providers.Track, error) {
	return nil, fmt.Errorf("spotify: not implemented yet")
}

func (p *Provider) SetVolume(level int) error {
	return fmt.Errorf("spotify: not implemented yet")
}
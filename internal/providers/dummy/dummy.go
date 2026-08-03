package dummy

import (
	"fmt"

	"github.com/tunahanyelmer/Tune/internal/providers"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Login() error {
	return nil
}

func (p *Provider) Play(query string) error {
	fmt.Println("▶ Playing:", query)
	return nil
}

func (p *Provider) Pause() error {
	fmt.Println("⏸ Paused")
	return nil
}

func (p *Provider) Next() error {
	fmt.Println("⏭ Next")
	return nil
}

func (p *Provider) Previous() error {
	fmt.Println("⏮ Previous")
	return nil
}

func (p *Provider) Current() (*providers.Track, error) {
	return &providers.Track{
		
		Title:  "Dummy Song",
		Artist: "Dummy Artist",
	}, nil
}

func (p *Provider) Search(query string) ([]providers.Track, error) {
	return []providers.Track{
		{
			
			Title:  query,
			Artist: "Dummy Artist",
		},
	}, nil
}

func (p *Provider) SetVolume(level int) error {
	fmt.Printf("Volume: %d\n", level)
	return nil
}
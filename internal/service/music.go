package service

import "github.com/tunahanyelmer/Tune/internal/providers"

type MusicService struct {
	provider providers.Provider
}

func NewMusicService(p providers.Provider) *MusicService {
	return &MusicService{
		provider: p,
	}
}

func (m *MusicService) Play(query string) error {
	return m.provider.Play(query)
}

func (m *MusicService) Pause() error {
	return m.provider.Pause()
}

func (m *MusicService) Next() error {
	return m.provider.Next()
}

func (m *MusicService) Previous() error {
	return m.provider.Previous()
}

func (m *MusicService) Login() error {
	return m.provider.Login()
}

func (m *MusicService) Current() (*providers.Track, error) {
	return m.provider.Current()
}

func (m *MusicService) Search(query string) ([]providers.Track, error) {
	return m.provider.Search(query)
}

func (m *MusicService) SetVolume(level int) error {
	return m.provider.SetVolume(level)
}
package registry

import (
	"fmt"

	"github.com/tunahanyelmer/Tune/internal/providers"
	"github.com/tunahanyelmer/Tune/internal/providers/dummy"
	"github.com/tunahanyelmer/Tune/internal/providers/spotify"
	"github.com/tunahanyelmer/Tune/internal/providers/youtube_music"
)

func Load(name string) (providers.Provider, error) {
	switch name {
	case "spotify":
		return spotify.NewSpotifyProvider(), nil
	case "youtube_music":
		return youtube_music.NewYouTubeMusicProvider(), nil
	case "dummy", "":
		return dummy.NewProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %q", name)
	}
}
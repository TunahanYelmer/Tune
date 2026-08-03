package daemon

import "github.com/tunahanyelmer/Tune/internal/providers"

type Action string

const (
	ActionLogin   Action = "login"
	ActionPlay    Action = "play"
	ActionPause   Action = "pause"
	ActionNext    Action = "next"
	ActionPrev    Action = "previous"
	ActionCurrent Action = "current"
	ActionSearch  Action = "search"
	ActionVolume  Action = "volume"
)

type Request struct {
	Action Action `json:"action"`
	Query  string `json:"query,omitempty"`
	Level  int    `json:"level,omitempty"`
}

type Response struct {
	OK     bool              `json:"ok"`
	Error  string            `json:"error,omitempty"`
	Track  *providers.Track  `json:"track,omitempty"`
	Tracks []providers.Track `json:"tracks,omitempty"`
}
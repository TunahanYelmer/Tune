package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tunahanyelmer/Tune/internal/daemon"
	"github.com/tunahanyelmer/Tune/internal/providers"
)

type mode int

const (
	modeNormal mode = iota
	modeSearch
)

type Model struct {
	provider  string
	current   *providers.Track
	results   []providers.Track
	cursor    int
	mode      mode
	searchBox string
	status    string
	loading   bool
	err       error
	quitting  bool
}

func New(provider string) Model {
	return Model{provider: provider}
}

// --- Bubble Tea lifecycle ---

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchCurrentCmd(), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tickMsg:
		return m, tea.Batch(fetchCurrentCmd(), tickCmd())

	case currentMsg:
		m.current = msg.track
		m.err = msg.err
		return m, nil

	case searchResultMsg:
		m.loading = false
		m.results = msg.tracks
		m.err = msg.err
		m.cursor = 0
		return m, nil

	case actionMsg:
		m.loading = false
		m.status = msg.status
		m.err = msg.err
		return m, fetchCurrentCmd()
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.mode == modeSearch {
		return m.handleSearchKey(msg)
	}
	return m.handleNormalKey(msg)
}

func (m Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit

	case " ":
		m.loading = true
		return m, sendActionCmd(daemon.ActionPause, "", 0, "⏯ Toggled play/pause")

	case "right":
		m.loading = true
		return m, sendActionCmd(daemon.ActionNext, "", 0, "⏭ Skipped to next")

	case "left":
		m.loading = true
		return m, sendActionCmd(daemon.ActionPrev, "", 0, "⏮ Went to previous")

	case "/":
		m.mode = modeSearch
		m.searchBox = ""
		return m, nil
	}
	return m, nil
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = modeNormal
		return m, nil

	case tea.KeyEnter:
		if len(m.results) > 0 && m.cursor < len(m.results) {
			track := m.results[m.cursor]
			m.mode = modeNormal
			m.loading = true
			query := track.Artist + " " + track.Title
			return m, sendActionCmd(daemon.ActionPlay, query, 0, "▶ Playing: "+track.Title)
		}
		m.loading = true
		query := m.searchBox
		return m, searchCmd(query)

	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case tea.KeyDown:
		if m.cursor < len(m.results)-1 {
			m.cursor++
		}
		return m, nil

	case tea.KeyBackspace:
		if len(m.searchBox) > 0 {
			m.searchBox = m.searchBox[:len(m.searchBox)-1]
		}
		return m, nil

	case tea.KeyRunes:
		m.searchBox += string(msg.Runes)
		return m, nil
	}
	return m, nil
}

// --- View ---

var (
	boxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Width(70)
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
)

func (m Model) View() string {
	if m.quitting {
		return "Bye 👋\n"
	}

	var b strings.Builder
	b.WriteString(renderLogo(m.provider) + "\n\n")

	if m.mode == modeSearch {
		b.WriteString("Search: " + m.searchBox + "\n\n")
		if len(m.results) == 0 && !m.loading {
			b.WriteString(dimStyle.Render("Type a query and press Enter to search") + "\n")
		}
		for i, t := range m.results {
			line := fmt.Sprintf("%s - %s", t.Artist, t.Title)
			if i == m.cursor {
				b.WriteString(selStyle.Render("▶ "+line) + "\n")
			} else {
				b.WriteString("  " + line + "\n")
			}
		}
		if m.loading {
			b.WriteString("\n" + dimStyle.Render("⏳ Searching...") + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("Enter search / Enter (on result) play / Esc cancel"))
	} else {
		if m.current != nil {
			b.WriteString(fmt.Sprintf("▶ %s\n", m.current.Title))
			b.WriteString(fmt.Sprintf("  %s\n", m.current.Artist))
			if m.current.Album != "" {
				b.WriteString(dimStyle.Render("  "+m.current.Album) + "\n")
			}
		} else {
			b.WriteString(dimStyle.Render("Nothing playing — press / to search") + "\n")
		}

		if m.loading {
			b.WriteString("\n" + dimStyle.Render("⏳ Working...") + "\n")
		} else if m.err != nil {
			b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("⚠ "+m.err.Error()) + "\n")
		} else if m.status != "" {
			b.WriteString("\n" + dimStyle.Render(m.status) + "\n")
		}

		b.WriteString("\n" + dimStyle.Render("Space Play/Pause   ← → Previous/Next   / Search   q Quit"))
	}

	return boxStyle.Render(b.String())
}

// --- messages & commands ---

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type currentMsg struct {
	track *providers.Track
	err   error
}

func fetchCurrentCmd() tea.Cmd {
	return func() tea.Msg {
		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionCurrent})
		if err != nil {
			return currentMsg{err: err}
		}
		if !resp.OK {
			return currentMsg{err: fmt.Errorf(resp.Error)}
		}
		return currentMsg{track: resp.Track}
	}
}

type searchResultMsg struct {
	tracks []providers.Track
	err    error
}

func searchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionSearch, Query: query})
		if err != nil {
			return searchResultMsg{err: err}
		}
		if !resp.OK {
			return searchResultMsg{err: fmt.Errorf(resp.Error)}
		}
		return searchResultMsg{tracks: resp.Tracks}
	}
}

type actionMsg struct {
	status string
	err    error
}

func sendActionCmd(action daemon.Action, query string, level int, successStatus string) tea.Cmd {
	return func() tea.Msg {
		resp, err := daemon.Send(daemon.Request{Action: action, Query: query, Level: level})
		if err != nil {
			return actionMsg{err: err}
		}
		if !resp.OK {
			return actionMsg{err: fmt.Errorf(resp.Error)}
		}
		return actionMsg{status: successStatus}
	}
}
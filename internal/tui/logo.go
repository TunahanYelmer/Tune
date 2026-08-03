package tui

import "github.com/charmbracelet/lipgloss"

// asciiLogos holds a small multi-line ASCII art block per provider,
// shown in the TUI header.
var asciiLogos = map[string]string{
	"spotify": `
 ███████╗██████╗  ██████╗ ████████╗██╗███████╗██╗   ██╗
 ██╔════╝██╔══██╗██╔═══██╗╚══██╔══╝██║██╔════╝╚██╗ ██╔╝
 ███████╗██████╔╝██║   ██║   ██║   ██║█████╗   ╚████╔╝ 
 ╚════██║██╔═══╝ ██║   ██║   ██║   ██║██╔══╝    ╚██╔╝  
 ███████║██║     ╚██████╔╝   ██║   ██║██║        ██║   
 ╚══════╝╚═╝      ╚═════╝    ╚═╝   ╚═╝╚═╝        ╚═╝   `,

	"youtube_music": `
 ▶ ██╗   ██╗████████╗    ███╗   ███╗██╗   ██╗███████╗██╗ ██████╗
   ╚██╗ ██╔╝╚══██╔══╝    ████╗ ████║██║   ██║██╔════╝██║██╔════╝
    ╚████╔╝    ██║       ██╔████╔██║██║   ██║███████╗██║██║     
     ╚██╔╝     ██║       ██║╚██╔╝██║██║   ██║╚════██║██║██║     
      ██║      ██║       ██║ ╚═╝ ██║╚██████╔╝███████║██║╚██████╗
      ╚═╝      ╚═╝       ╚═╝     ╚═╝ ╚═════╝ ╚══════╝╚═╝ ╚═════╝`,

	"dummy": `
 ██████╗ ██╗   ██╗███╗   ███╗███╗   ███╗██╗   ██╗
 ██╔══██╗██║   ██║████╗ ████║████╗ ████║╚██╗ ██╔╝
 ██║  ██║██║   ██║██╔████╔██║██╔████╔██║ ╚████╔╝ 
 ██║  ██║██║   ██║██║╚██╔╝██║██║╚██╔╝██║  ╚██╔╝  
 ██████╔╝╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║   ██║   
 ╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝   ╚═╝   `,
}

var logoColors = map[string]string{
	"spotify":       "#1DB954", // Spotify green
	"youtube_music": "#FF0000", // YouTube red
	"dummy":         "#888888",
}

func renderLogo(provider string) string {
	art, ok := asciiLogos[provider]
	if !ok {
		art = "\n  ♪ TUNE ♪\n"
	}
	color, ok := logoColors[provider]
	if !ok {
		color = "#888888"
	}

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))
	return style.Render(art)
}
package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#8A2BE2") // Vibrant purple
	ColorSecondary = lipgloss.Color("#00D7FF") // Electric blue
	ColorSuccess   = lipgloss.Color("#00FF66") // Neon green
	ColorWarning   = lipgloss.Color("#FFB300") // Yellow/Amber
	ColorDanger    = lipgloss.Color("#FF3366") // Neon Red/Pink
	ColorMuted     = lipgloss.Color("#6272A4") // Dracula gray
	ColorDarkBg    = lipgloss.Color("#1E1E2E") // Catppuccin mocha bg
	ColorHighlight = lipgloss.Color("#F8F8F2") // Off-white

	// Base Styles
	StyleTitle = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorHighlight).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Italic(true).
			MarginBottom(1)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorMuted).
			PaddingTop(1)

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	StyleSelected = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorHighlight).
			Bold(true)

	StyleUnselected = lipgloss.NewStyle().
			Foreground(ColorHighlight)

	// Status Styles
	StyleSolved = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleUnsolved = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleMissed = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	// Difficulty Styles
	StyleEasy = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StyleMedium = lipgloss.NewStyle().
			Foreground(ColorWarning)

	StyleHard = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	// Category Badges
	StyleCategory = lipgloss.NewStyle().
			Background(lipgloss.Color("#282A36")).
			Foreground(ColorSecondary).
			Padding(0, 1).
			Bold(true)

	// Table Headers
	StyleHeader = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true).
			Underline(true)

	// Sidebar Details Panel
	StyleDetailTitle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Bold(true).
				MarginBottom(1)

	StyleDetailLabel = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Width(12)

	StyleDetailVal = lipgloss.NewStyle().
			Foreground(ColorHighlight)
)

package tui

import (
	"fmt"
	"strings"

	"github.com/WissemJderi/go-ctf-tracker/pkg/storage"
	"github.com/charmbracelet/lipgloss"
)

const (
	DetailBoxWidth      = 75
	DetailBoxPadding    = 6
	ColumnCTFName       = 20
	ColumnChallengeName = 25
	ColumnCategory      = 10
	ColumnPoints        = 6
	ColumnDiff          = 6
	ColumnStatus        = 10
	ColumnHard          = 6
)

func (m Model) View() string {
	var s strings.Builder

	// 1. Title/Header Banner
	s.WriteString(StyleTitle.Render("⚡ CTF CHALLENGE TRACKER ⚡") + "\n")
	s.WriteString(StyleSubtitle.Render("Manage live challenges, track missed ones, link writeups!") + "\n\n")

	// 2. Info / Error Messages
	if m.errorMsg != "" {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorDanger).Bold(true).Render("✖ "+m.errorMsg) + "\n\n")
	} else if m.infoMsg != "" {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("✔ "+m.infoMsg) + "\n\n")
	}

	// 3. Render State View
	switch m.state {
	case stateList:
		s.WriteString(m.listView())
	case stateDetail:
		s.WriteString(m.detailView())
	case stateAdd:
		s.WriteString(StyleDetailTitle.Render("➕ ADD NEW CHALLENGE") + "\n")
		s.WriteString(m.form.View())
	case stateEdit:
		s.WriteString(StyleDetailTitle.Render("📝 EDIT CHALLENGE") + "\n")
		s.WriteString(m.form.View())
	case stateFilterCTFPrompt:
		s.WriteString(StyleDetailTitle.Render("🔍 FILTER BY CTF NAME") + "\n")
		s.WriteString(m.ctfInput.View())
	case stateDeleteConfirm:
		s.WriteString(m.deleteConfirmView())
	}

	return s.String()
}

func (m Model) listView() string {
	var s strings.Builder

	// Render Active Filters
	var activeFilters []string
	if m.filters.CTF != "" {
		activeFilters = append(activeFilters, fmt.Sprintf("CTF: %s", m.filters.CTF))
	}
	if m.filters.Status != "" {
		activeFilters = append(activeFilters, fmt.Sprintf("Status: %s", m.filters.Status))
	}
	if m.filters.FlaggedHard {
		activeFilters = append(activeFilters, "Hard Only 🔥")
	}
	if m.filters.MissedOnly {
		activeFilters = append(activeFilters, "Missed Only ✖")
	}

	if len(activeFilters) > 0 {
		filterBar := lipgloss.NewStyle().
			Background(lipgloss.Color("#44475A")).
			Foreground(ColorWarning).
			Padding(0, 1).
			Bold(true).
			Render("⏳ Filters: [ " + strings.Join(activeFilters, " | ") + " ] (Press 'x' to clear)")
		s.WriteString(filterBar + "\n\n")
	}

	if len(m.filtered) == 0 {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render("   No challenges found. Press 'a' to add a new challenge!"))
		s.WriteString("\n\n" + m.listHelpView())
		return s.String()
	}

	// Print table header
	headerLine := fmt.Sprintf("  %s  %s  %s  %s  %s  %s  %s  %s",
		lipgloss.NewStyle().Width(4).Render("ID"),
		lipgloss.NewStyle().Width(20).Render("CTF NAME"),
		lipgloss.NewStyle().Width(25).Render("CHALLENGE NAME"),
		lipgloss.NewStyle().Width(10).Render("CATEGORY"),
		lipgloss.NewStyle().Width(6).Render("POINTS"),
		lipgloss.NewStyle().Width(6).Render("DIFF"),
		lipgloss.NewStyle().Width(10).Render("STATUS"),
		lipgloss.NewStyle().Width(6).Render("HARD?"),
	)
	s.WriteString(StyleHeader.Render(headerLine) + "\n")

	for i, ch := range m.filtered {
		isCursor := i == m.cursor

		// Render ID
		idStr := ch.ID
		if len(idStr) > 4 {
			idStr = idStr[:4]
		}

		// Render CTF & Challenge Names with truncation to prevent wrapping
		ctfStr := truncateString(ch.CTFName, 20)
		nameStr := truncateString(ch.Name, 25)

		// Category
		catStr := ch.Category

		// Points
		ptsStr := fmt.Sprintf("%d", ch.Points)

		// Difficulty styling
		var diffStyle lipgloss.Style
		switch strings.ToLower(ch.Difficulty) {
		case "easy":
			diffStyle = StyleEasy
		case "medium":
			diffStyle = StyleMedium
		case "hard":
			diffStyle = StyleHard
		default:
			diffStyle = StyleUnselected
		}

		// Status styling & text
		var statusText string
		var statusStyle lipgloss.Style
		switch ch.Status {
		case storage.StatusSolved:
			statusText = "✔ Solved"
			statusStyle = StyleSolved
		case storage.StatusMissed:
			statusText = "✖ Missed"
			statusStyle = StyleMissed
		default:
			statusText = "◦ Unsolved"
			statusStyle = StyleUnsolved
		}

		// Hard text & style
		var hardText string
		var hardStyle lipgloss.Style
		if ch.FlaggedHard {
			hardText = "🔥 YES"
			hardStyle = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
		} else {
			hardText = "no"
			hardStyle = lipgloss.NewStyle().Foreground(ColorMuted)
		}

		cellID := lipgloss.NewStyle().Width(4).Render(idStr)
		cellCTF := lipgloss.NewStyle().Width(20).Render(ctfStr)
		cellName := lipgloss.NewStyle().Width(25).Render(nameStr)
		cellCat := lipgloss.NewStyle().Width(10).Render(catStr)
		cellPts := lipgloss.NewStyle().Width(6).Render(ptsStr)
		cellDiff := diffStyle.Width(6).Render(ch.Difficulty)
		cellStatus := statusStyle.Width(10).Render(statusText)
		cellHard := hardStyle.Width(6).Render(hardText)

		rowContent := fmt.Sprintf("%s  %s  %s  %s  %s  %s  %s  %s",
			cellID, cellCTF, cellName, cellCat, cellPts, cellDiff, cellStatus, cellHard)

		if isCursor {
			s.WriteString(StyleSelected.Render("▸ "+rowContent) + "\n")
		} else {
			s.WriteString(StyleUnselected.Render("  "+rowContent) + "\n")
		}
	}

	s.WriteString("\n" + m.listHelpView())
	return s.String()
}

func (m Model) detailView() string {
	ch, ok := m.getSelectedChallenge()
	if !ok {
		return "No challenge selected."
	}

	var s strings.Builder

	border := strings.Repeat("━", DetailBoxWidth)

	s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┏"+border+"┓") + "\n")

	// Title
	titleLine := fmt.Sprintf(" CHALLENGE DETAILS: %s ", ch.Name)
	titlePad := (DetailBoxWidth - len(titleLine)) / 2
	if titlePad < 0 {
		titlePad = 0
	}
	s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃") +
		strings.Repeat(" ", titlePad) +
		lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Render(titleLine) +
		strings.Repeat(" ", DetailBoxWidth-titlePad-len(titleLine)) +
		lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃") + "\n")

	s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┣"+border+"┫") + "\n")

	// Helper for printing key-value rows in the box
	printRow := func(label, val string, valStyle lipgloss.Style) {
		lbl := StyleDetailLabel.Render(label + ":")
		v := valStyle.Render(val)
		contentLen := lipgloss.Width(lbl) + lipgloss.Width(v)
		pad := DetailBoxWidth - contentLen - 4
		if pad < 0 {
			pad = 0
		}
		s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃ ") +
			lbl + " " + v +
			strings.Repeat(" ", pad) + " " +
			lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃") + "\n")
	}

	// Fields
	printRow("ID", ch.ID, StyleDetailVal)
	printRow("CTF Name", ch.CTFName, StyleDetailVal)
	printRow("Category", ch.Category, StyleCategory)
	printRow("Points", fmt.Sprintf("%d", ch.Points), StyleDetailVal)

	var diffStyle lipgloss.Style
	switch strings.ToLower(ch.Difficulty) {
	case "easy":
		diffStyle = StyleEasy
	case "medium":
		diffStyle = StyleMedium
	case "hard":
		diffStyle = StyleHard
	}
	printRow("Difficulty", ch.Difficulty, diffStyle)

	var statStyle lipgloss.Style
	switch ch.Status {
	case storage.StatusSolved:
		statStyle = StyleSolved
	case storage.StatusMissed:
		statStyle = StyleMissed
	default:
		statStyle = StyleUnsolved
	}
	printRow("Status", string(ch.Status), statStyle)

	var hardVal string
	var hardStyle lipgloss.Style
	if ch.FlaggedHard {
		hardVal = "YES (🔥 LIVE HARD CHALLENGE)"
		hardStyle = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
	} else {
		hardVal = "NO"
		hardStyle = StyleUnselected
	}
	printRow("Flagged Hard", hardVal, hardStyle)

	s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┣"+border+"┫") + "\n")

	// Notes & Writeup Sections
	printSection := func(title, content string, isLink bool) {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃ ") +
			lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Render(title+":") +
			strings.Repeat(" ", DetailBoxWidth-len(title)-4) + " " +
			lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃") + "\n")

		if content == "" {
			content = "None."
		}

		// Wrap and indent content lines
		lines := wordWrap(content, DetailBoxWidth-DetailBoxPadding)
		for _, line := range lines {
			var lineStyle lipgloss.Style
			if isLink && content != "None." {
				lineStyle = lipgloss.NewStyle().Foreground(ColorSecondary).Underline(true)
			} else {
				lineStyle = StyleDetailVal
			}
			renderedLine := lineStyle.Render(line)
			pad := DetailBoxWidth - lipgloss.Width(renderedLine) - 4
			if pad < 0 {
				pad = 0
			}
			s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃   ") +
				renderedLine +
				strings.Repeat(" ", pad) + " " +
				lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃") + "\n")
		}
	}

	printSection("WRITEUP LINK", ch.WriteupURL, true)
	s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┃"+strings.Repeat(" ", DetailBoxWidth)+"┃") + "\n")
	printSection("NOTES / CLUES", ch.Notes, false)

	s.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("┗"+border+"┛") + "\n\n")

	// Custom detail help controls
	helpLines := []string{
		"Keys: [ESC/v/Enter] list | [e] edit | [s] cycle status | [h] toggle hard flag",
	}
	for _, l := range helpLines {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render(l) + "\n")
	}

	return s.String()
}

func (m Model) deleteConfirmView() string {
	ch, ok := m.getSelectedChallenge()
	if !ok {
		return ""
	}

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Background(ColorDanger).Foreground(ColorHighlight).Bold(true).Padding(0, 2).Render("⚠ DANGER: DELETE CHALLENGE ⚠") + "\n\n")
	s.WriteString(fmt.Sprintf("Are you absolutely sure you want to delete challenge %s?\n\n", lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Render(ch.Name)))
	s.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("[y] Yes, Delete") + "   " + lipgloss.NewStyle().Foreground(ColorMuted).Render("[n/ESC] No, Cancel") + "\n")

	return s.String()
}

func (m Model) listHelpView() string {
	var s strings.Builder
	s.WriteString(StyleHelp.Render("─── CONTROLS ───────────────────────────────────────────────────────────────────────────") + "\n")

	helpGrid := [][]string{
		{"a", "add challenge", "s", "cycle status (unsolved->solved->missed)"},
		{"e", "edit challenge", "h", "toggle hard flag (🔥)"},
		{"d", "delete challenge", "t", "filter by CTF name"},
		{"v/enter", "view challenge details", "c", "cycle status filter"},
		{"f", "toggle flagged hard filter", "m", "toggle missed filter"},
		{"x", "clear all filters", "q", "quit tracker"},
	}

	for _, row := range helpGrid {
		key1 := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Width(8).Render(row[0])
		desc1 := lipgloss.NewStyle().Foreground(ColorHighlight).Width(25).Render(row[1])
		key2 := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Width(8).Render(row[2])
		desc2 := lipgloss.NewStyle().Foreground(ColorHighlight).Render(row[3])
		s.WriteString(fmt.Sprintf("  %s %s │ %s %s\n", key1, desc1, key2, desc2))
	}

	return s.String()
}

// Formatting helpers
func truncateString(str string, maxLen int) string {
	if maxLen < 0 {
		maxLen = 0
	}
	if len(str) > maxLen {
		if maxLen > 3 {
			return str[:maxLen-3] + "..."
		}
		return str[:maxLen]
	}
	return str
}

func wordWrap(text string, maxLen int) []string {
	var lines []string
	paragraphs := strings.Split(text, "\n")

	for _, p := range paragraphs {
		if p == "" {
			lines = append(lines, "")
			continue
		}

		words := strings.Fields(p)
		if len(words) == 0 {
			continue
		}

		currentLine := words[0]
		for _, word := range words[1:] {
			if len(currentLine)+1+len(word) <= maxLen {
				currentLine += " " + word
			} else {
				lines = append(lines, currentLine)
				currentLine = word
			}
		}
		lines = append(lines, currentLine)
	}
	return lines
}

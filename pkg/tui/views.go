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
	header := m.headerView()
	controls := m.controlsView()

	headerH := renderedHeight(header)
	controlsH := renderedHeight(controls)
	// One line is reserved for the blank separator between header and body.
	bodyH := m.height - headerH - controlsH - 1
	if bodyH < 0 {
		bodyH = 0
	}
	body := strings.TrimRight(m.bodyView(bodyH), "\n")

	var s strings.Builder
	s.WriteString(header)
	if body != "" {
		s.WriteString("\n\n" + body)
	}
	used := renderedHeight(s.String())
	// The controls block reuses the bottom-most padding row as its own top
	// margin, so it consumes controlsH-1 rows of new space.
	if gap := m.height - used - controlsH + 1; gap > 0 {
		s.WriteString(strings.Repeat("\n", gap))
	}
	s.WriteString(controls)

	return s.String()
}

func (m Model) headerView() string {
	var h strings.Builder
	h.WriteString(StyleTitle.MarginBottom(0).Render("⚡ CTF CHALLENGE TRACKER ⚡"))
	h.WriteString("\n\n")
	h.WriteString(StyleSubtitle.MarginBottom(0).Render("Manage live challenges, track missed ones, link writeups!"))
	if m.errorMsg != "" {
		h.WriteString("\n\n" + lipgloss.NewStyle().Foreground(ColorDanger).Bold(true).Render("✖ "+m.errorMsg))
	} else if m.infoMsg != "" {
		h.WriteString("\n\n" + lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("✔ "+m.infoMsg))
	}
	return h.String()
}

// renderedHeight counts visible lines a rendered block occupies, ignoring any
// trailing newlines.
func renderedHeight(s string) int {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func (m Model) controlsView() string {
	switch m.state {
	case stateList:
		return m.listHelpView()
	case stateDetail:
		return StyleHelp.Render("Keys: [ESC/v/Enter] list | [e] edit | [s] cycle status | [h] toggle hard flag")
	case stateAdd, stateEdit:
		return StyleHelp.Render("Keys: [Tab] next field | [Enter] submit | [Esc] cancel")
	case stateFilterCTFPrompt:
		return StyleHelp.Render("Keys: [Enter] apply filter | [Esc] cancel")
	case stateDeleteConfirm:
		return StyleHelp.Render("Keys: [y] confirm delete | [n/Esc] cancel")
	}
	return ""
}

func (m Model) bodyView(maxBodyHeight int) string {
	switch m.state {
	case stateList:
		return m.listBodyView(maxBodyHeight)
	case stateDetail:
		return m.detailBodyView()
	case stateAdd:
		return StyleDetailTitle.Render("➕ ADD NEW CHALLENGE") + "\n" + m.form.View()
	case stateEdit:
		return StyleDetailTitle.Render("📝 EDIT CHALLENGE") + "\n" + m.form.View()
	case stateFilterCTFPrompt:
		return StyleDetailTitle.Render("🔍 FILTER BY CTF NAME") + "\n" + m.ctfInput.View()
	case stateDeleteConfirm:
		return m.deleteConfirmView()
	}
	return ""
}

func (m Model) filterBarView() string {
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

	if len(activeFilters) == 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#44475A")).
		Foreground(ColorWarning).
		Padding(0, 1).
		Bold(true).
		Render("⏳ Filters: [ " + strings.Join(activeFilters, " | ") + " ] (Press 'x' to clear)")
}

func (m Model) listBodyView(maxBodyHeight int) string {
	if len(m.filtered) == 0 {
		return lipgloss.NewStyle().Foreground(ColorMuted).Render("   No challenges found. Press 'a' to add a new challenge!")
	}

	var lead strings.Builder
	if bar := m.filterBarView(); bar != "" {
		lead.WriteString(bar + "\n\n")
	}

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
	lead.WriteString(StyleHeader.Render(headerLine) + "\n")
	leadLines := renderedHeight(lead.String())

	// Reserve 2 lines for the "more above / below" scroll indicators.
	rowsAvail := maxBodyHeight - leadLines - 2
	if rowsAvail < 1 {
		rowsAvail = 1
	}
	if rowsAvail > len(m.filtered) {
		rowsAvail = len(m.filtered)
	}

	start, end := m.listWindow(rowsAvail)

	var s strings.Builder
	s.WriteString(lead.String())

	if start > 0 {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("  ↕ ... %d more above", start)) + "\n")
	}
	for i := start; i < end; i++ {
		s.WriteString(m.rowView(i) + "\n")
	}
	if end < len(m.filtered) {
		s.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("  ↕ ... %d more below", len(m.filtered)-end)) + "\n")
	}

	return s.String()
}

// listWindow returns the [start, end) slice of m.filtered that fits in
// rowsAvail lines while keeping m.cursor visible.
func (m Model) listWindow(rowsAvail int) (int, int) {
	total := len(m.filtered)
	if rowsAvail >= total {
		return 0, total
	}

	start := 0
	end := rowsAvail
	if m.cursor >= end {
		start = m.cursor - rowsAvail + 1
		end = m.cursor + 1
		if end > total {
			end = total
			start = end - rowsAvail
		}
	}
	if start < 0 {
		start = 0
		end = start + rowsAvail
	}
	return start, end
}

func (m Model) rowView(i int) string {
	ch := m.filtered[i]
	isCursor := i == m.cursor

	idStr := ch.ID
	if len(idStr) > 4 {
		idStr = idStr[:4]
	}

	// Truncate names to prevent wrapping
	ctfStr := truncateString(ch.CTFName, 20)
	nameStr := truncateString(ch.Name, 25)

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

	var rowStyle lipgloss.Style
	if isCursor {
		rowStyle = StyleSelected
	} else {
		rowStyle = StyleUnselected
	}

	cellID := rowStyle.Width(4).Render(idStr)
	cellCTF := rowStyle.Width(20).Render(ctfStr)
	cellName := rowStyle.Width(25).Render(nameStr)
	cellCat := rowStyle.Width(10).Render(ch.Category)
	cellPts := rowStyle.Width(6).Render(fmt.Sprintf("%d", ch.Points))
	cellDiff := diffStyle.Width(6).Render(ch.Difficulty)
	cellStatus := statusStyle.Width(10).Render(statusText)
	cellHard := hardStyle.Width(6).Render(hardText)

	rowContent := fmt.Sprintf("%s  %s  %s  %s  %s  %s  %s  %s",
		cellID, cellCTF, cellName, cellCat, cellPts, cellDiff, cellStatus, cellHard)

	return rowStyle.Render("▸ " + rowContent)
}

func (m Model) detailBodyView() string {
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
	s.WriteString(StyleHelp.Render("─── CONTROLS ───────────────────────────────────────────────────────────────") + "\n")

	helpGrid := [][]string{
		{"a", "add challenge", "s", "cycle status (unsolved->solved->missed)"},
		{"e", "edit challenge", "h", "toggle hard flag (🔥)"},
		{"d", "delete challenge", "t", "filter by CTF name"},
		{"v/enter", "view challenge details", "c", "cycle status filter"},
		{"f", "toggle hard only", "m", "toggle missed filter"},
		{"x", "clear all filters", "q", "quit tracker"},
	}

	for i, row := range helpGrid {
		key1 := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Width(8).Render(row[0])
		desc1 := lipgloss.NewStyle().Foreground(ColorHighlight).Width(25).Render(row[1])
		key2 := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Width(8).Render(row[2])
		desc2 := lipgloss.NewStyle().Foreground(ColorHighlight).Render(row[3])
		s.WriteString(fmt.Sprintf("  %s %s │ %s %s", key1, desc1, key2, desc2))
		if i < len(helpGrid)-1 {
			s.WriteString("\n")
		}
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

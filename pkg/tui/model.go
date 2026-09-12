package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/WissemJderi/go-ctf-tracker/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	stateList sessionState = iota
	stateDetail
	stateAdd
	stateEdit
	stateDeleteConfirm
	stateFilterCTFPrompt
)

type FilterSettings struct {
	CTF          string
	Category     string
	Status       string // "", "Unsolved", "Solved", "Missed"
	FlaggedHard  bool
	MissedOnly   bool
}

type Model struct {
	storage       *storage.JSONStorage
	challenges    []storage.Challenge
	filtered      []storage.Challenge
	cursor        int
	state         sessionState
	filters       FilterSettings
	form          FormModel
	ctfInput      FormModel // reused simple form for entering CTF name to filter
	selectedID    string
	width         int
	height        int
	errorMsg      string
	infoMsg       string
}

func NewModel(store *storage.JSONStorage) Model {
	return Model{
		storage: store,
		state:   stateList,
	}
}

// Commands
func (m Model) loadChallenges() tea.Cmd {
	return func() tea.Msg {
		challs, err := m.storage.Load()
		if err != nil {
			return errMsg{err}
		}
		return challengesLoadedMsg{challs}
	}
}

func (m Model) saveChallenge(ch storage.Challenge, isNew bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if isNew {
			err = m.storage.Add(ch)
		} else {
			err = m.storage.Update(ch)
		}
		if err != nil {
			return errMsg{err}
		}
		return operationSuccessMsg(fmt.Sprintf("Challenge successfully %s!", map[bool]string{true: "added", false: "updated"}[isNew]))
	}
}

func (m Model) deleteChallengeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.storage.Delete(id)
		if err != nil {
			return errMsg{err}
		}
		return operationSuccessMsg("Challenge successfully deleted!")
	}
}

// Messages
type errMsg struct{ err error }
type challengesLoadedMsg struct{ challenges []storage.Challenge }
type operationSuccessMsg string

func (m Model) Init() tea.Cmd {
	return m.loadChallenges()
}

func (m Model) applyFilters() []storage.Challenge {
	var filtered []storage.Challenge
	for _, c := range m.challenges {
		// Filter by CTF Name
		if m.filters.CTF != "" && !strings.Contains(strings.ToLower(c.CTFName), strings.ToLower(m.filters.CTF)) {
			continue
		}
		// Filter by Status
		if m.filters.Status != "" && string(c.Status) != m.filters.Status {
			continue
		}
		// Filter by Flagged Hard
		if m.filters.FlaggedHard && !c.FlaggedHard {
			continue
		}
		// Filter by Missed Only
		if m.filters.MissedOnly && c.Status != storage.StatusMissed {
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}

func normalizeCursor(challenges, filtered []storage.Challenge, cursor int) int {
	if len(filtered) == 0 {
		return 0
	}
	if cursor >= len(filtered) {
		return len(filtered) - 1
	}
	if cursor < 0 {
		return 0
	}
	return cursor
}

func rotateStatus(status storage.ChallengeStatus) storage.ChallengeStatus {
	switch status {
	case storage.StatusUnsolved:
		return storage.StatusSolved
	case storage.StatusSolved:
		return storage.StatusMissed
	case storage.StatusMissed:
		return storage.StatusUnsolved
	default:
		return storage.StatusUnsolved
	}
}

func (m Model) getSelectedChallenge() (storage.Challenge, bool) {
	if len(m.filtered) == 0 || m.cursor < 0 || m.cursor >= len(m.filtered) {
		return storage.Challenge{}, false
	}
	return m.filtered[m.cursor], true
}

// mutateSelected applies fn to the currently selected challenge and writes the
// result back into m.challenges by ID (not by cursor position, since m.cursor
// indexes m.filtered, which can differ from m.challenges once a filter is active).
func (m *Model) mutateSelected(fn func(*storage.Challenge)) (storage.Challenge, bool) {
	ch, ok := m.getSelectedChallenge()
	if !ok {
		return storage.Challenge{}, false
	}
	fn(&ch)
	for i := range m.challenges {
		if m.challenges[i].ID == ch.ID {
			m.challenges[i] = ch
			break
		}
	}
	return ch, true
}

// refilter recomputes the filtered list and clamps the cursor to it. Call
// after changing m.challenges or m.filters.
func (m *Model) refilter() {
	m.filtered = m.applyFilters()
	m.cursor = normalizeCursor(m.challenges, m.filtered, m.cursor)
}

func (m Model) initAddForm() FormModel {
	fields := []FormField{
		{Key: "ctf_name", Label: "CTF Name", Type: FieldText, Placeholder: "e.g., HTB Cyber Apocalypse"},
		{Key: "name", Label: "Challenge Name", Type: FieldText, Placeholder: "e.g., Simple Encryptor"},
		{Key: "category", Label: "Category", Type: FieldSelect, Options: []string{"Crypto", "Pwn", "Web", "Rev", "Misc", "Forensics", "OSINT", "Cloud"}},
		{Key: "points", Label: "Points", Type: FieldText, Placeholder: "e.g., 100"},
		{Key: "difficulty", Label: "Difficulty", Type: FieldSelect, Options: []string{"Easy", "Medium", "Hard"}},
		{Key: "flagged_hard", Label: "Flagged Hard?", Type: FieldBool},
		{Key: "writeup_url", Label: "Writeup URL", Type: FieldText, Placeholder: "e.g., https://github.com/.../writeup"},
		{Key: "notes", Label: "Notes / Clues", Type: FieldTextArea, Placeholder: "Enter any hints, ideas, or observations..."},
	}

	form := NewFormModel(fields, func(textVals map[string]string, boolVals map[string]bool) tea.Cmd {
		pts, err := strconv.Atoi(textVals["points"])
		if err != nil {
			return func() tea.Msg {
				return errMsg{fmt.Errorf("points must be a valid integer: %w", err)}
			}
		}
		ch := storage.Challenge{
			CTFName:     textVals["ctf_name"],
			Name:        textVals["name"],
			Category:    textVals["category"],
			Points:      pts,
			Difficulty:  textVals["difficulty"],
			Status:      storage.StatusUnsolved,
			FlaggedHard: boolVals["flagged_hard"],
			WriteupURL:  textVals["writeup_url"],
			Notes:       textVals["notes"],
		}
		return m.saveChallenge(ch, true)
	}, func() tea.Cmd {
		return func() tea.Msg { return cancelFormMsg{} }
	})

	return form
}

func (m Model) initEditForm(ch storage.Challenge) FormModel {
	fields := []FormField{
		{Key: "ctf_name", Label: "CTF Name", Type: FieldText, Placeholder: "e.g., HTB Cyber Apocalypse"},
		{Key: "name", Label: "Challenge Name", Type: FieldText, Placeholder: "e.g., Simple Encryptor"},
		{Key: "category", Label: "Category", Type: FieldSelect, Options: []string{"Crypto", "Pwn", "Web", "Rev", "Misc", "Forensics", "OSINT", "Cloud"}},
		{Key: "points", Label: "Points", Type: FieldText, Placeholder: "e.g., 100"},
		{Key: "difficulty", Label: "Difficulty", Type: FieldSelect, Options: []string{"Easy", "Medium", "Hard"}},
		{Key: "flagged_hard", Label: "Flagged Hard?", Type: FieldBool},
		{Key: "writeup_url", Label: "Writeup URL", Type: FieldText, Placeholder: "e.g., https://github.com/.../writeup"},
		{Key: "notes", Label: "Notes / Clues", Type: FieldTextArea, Placeholder: "Enter any hints, ideas, or observations..."},
	}

	form := NewFormModel(fields, func(textVals map[string]string, boolVals map[string]bool) tea.Cmd {
		pts, err := strconv.Atoi(textVals["points"])
		if err != nil {
			return func() tea.Msg {
				return errMsg{fmt.Errorf("points must be a valid integer: %w", err)}
			}
		}
		ch.CTFName = textVals["ctf_name"]
		ch.Name = textVals["name"]
		ch.Category = textVals["category"]
		ch.Points = pts
		ch.Difficulty = textVals["difficulty"]
		ch.FlaggedHard = boolVals["flagged_hard"]
		ch.WriteupURL = textVals["writeup_url"]
		ch.Notes = textVals["notes"]
		return m.saveChallenge(ch, false)
	}, func() tea.Cmd {
		return func() tea.Msg { return cancelFormMsg{} }
	})

	textVals := map[string]string{
		"ctf_name":    ch.CTFName,
		"name":        ch.Name,
		"category":    ch.Category,
		"points":      strconv.Itoa(ch.Points),
		"difficulty":  ch.Difficulty,
		"writeup_url": ch.WriteupURL,
		"notes":       ch.Notes,
	}
	boolVals := map[string]bool{
		"flagged_hard": ch.FlaggedHard,
	}
	form.SetValues(textVals, boolVals)
	return form
}

func (m Model) initCTFFilterPrompt() FormModel {
	fields := []FormField{
		{Key: "ctf_filter", Label: "CTF Name Filter", Type: FieldText, Placeholder: "Enter CTF name, e.g., HackTheBox (ESC/Enter)"},
	}
	return NewFormModel(fields, func(textVals map[string]string, boolVals map[string]bool) tea.Cmd {
		return func() tea.Msg { return applyCTFFilterMsg{textVals["ctf_filter"]} }
	}, func() tea.Cmd {
		return func() tea.Msg { return cancelFormMsg{} }
	})
}

type cancelFormMsg struct{}
type applyCTFFilterMsg struct{ ctf string }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case errMsg:
		m.errorMsg = msg.err.Error()
		m.infoMsg = ""
		return m, nil

	case operationSuccessMsg:
		m.infoMsg = string(msg)
		m.errorMsg = ""
		m.state = stateList
		return m, m.loadChallenges()

	case challengesLoadedMsg:
		m.challenges = msg.challenges
		m.refilter()
		return m, nil

	case cancelFormMsg:
		m.state = stateList
		return m, nil

	case applyCTFFilterMsg:
		m.filters.CTF = msg.ctf
		m.state = stateList
		m.refilter()
		return m, nil
	}

	// Route based on state
	switch m.state {
	case stateAdd, stateEdit:
		var formCmd tea.Cmd
		m.form, formCmd = m.form.Update(msg)
		return m, formCmd

	case stateFilterCTFPrompt:
		var formCmd tea.Cmd
		m.ctfInput, formCmd = m.ctfInput.Update(msg)
		return m, formCmd

	case stateDeleteConfirm:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				ch, ok := m.getSelectedChallenge()
				if ok {
					m.state = stateList
					return m, m.deleteChallengeCmd(ch.ID)
				}
				m.state = stateList
			case "n", "N", "esc":
				m.state = stateList
			}
		}
		return m, nil

	case stateDetail:
		ch, ok := m.getSelectedChallenge()
		if !ok {
			m.state = stateList
			return m, nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q", "v", "enter":
				m.state = stateList
			case "e":
				m.form = m.initEditForm(ch)
				m.state = stateEdit
			case "h":
				ch, _ := m.mutateSelected(func(c *storage.Challenge) { c.FlaggedHard = !c.FlaggedHard })
				return m, m.saveChallenge(ch, false)
			case "s":
				ch, _ := m.mutateSelected(func(c *storage.Challenge) { c.Status = rotateStatus(c.Status) })
				return m, m.saveChallenge(ch, false)
			}
		}
		return m, nil

	case stateList:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}

			case "a":
				m.form = m.initAddForm()
				m.state = stateAdd

			case "e":
				ch, ok := m.getSelectedChallenge()
				if ok {
					m.form = m.initEditForm(ch)
					m.state = stateEdit
				}

			case "d":
				_, ok := m.getSelectedChallenge()
				if ok {
					m.state = stateDeleteConfirm
				}

			case "v", "enter":
				_, ok := m.getSelectedChallenge()
				if ok {
					m.state = stateDetail
				}

			case "h": // Toggle Flagged Hard instantly
				ch, ok := m.mutateSelected(func(c *storage.Challenge) { c.FlaggedHard = !c.FlaggedHard })
				if ok {
					return m, m.saveChallenge(ch, false)
				}

			case "s": // Cycle Status instantly (Unsolved -> Solved -> Missed -> Unsolved)
				ch, ok := m.mutateSelected(func(c *storage.Challenge) { c.Status = rotateStatus(c.Status) })
				if ok {
					return m, m.saveChallenge(ch, false)
				}

			// FILTERS KEYBINDINGS
			case "f": // Toggle show only Flagged Hard
				m.filters.FlaggedHard = !m.filters.FlaggedHard
				m.refilter()

			case "m": // Toggle show only Missed
				m.filters.MissedOnly = !m.filters.MissedOnly
				m.refilter()

			case "t": // Filter by CTF prompt
				m.ctfInput = m.initCTFFilterPrompt()
				m.state = stateFilterCTFPrompt

			case "c": // Cycle Status Filter (All -> Unsolved -> Solved -> Missed -> All)
				switch m.filters.Status {
				case "":
					m.filters.Status = string(storage.StatusUnsolved)
				case string(storage.StatusUnsolved):
					m.filters.Status = string(storage.StatusSolved)
				case string(storage.StatusSolved):
					m.filters.Status = string(storage.StatusMissed)
				case string(storage.StatusMissed):
					m.filters.Status = ""
				}
				m.refilter()

			case "x": // Clear all filters
				m.filters = FilterSettings{}
				m.refilter()
				m.infoMsg = "Filters cleared!"
			}
		}
	}

	return m, cmd
}

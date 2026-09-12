package tui

import (
	"testing"

	"github.com/WissemJderi/go-ctf-tracker/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func keyMsg(s string) tea.KeyMsg {
	switch s {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

func newTestModel(challenges []storage.Challenge) Model {
	m := Model{
		state:      stateList,
		challenges: challenges,
	}
	m.filtered = m.applyFilters()
	return m
}

func TestListNavigation(t *testing.T) {
	m := newTestModel([]storage.Challenge{{ID: "1"}, {ID: "2"}, {ID: "3"}})

	mi, _ := m.Update(keyMsg("down"))
	m = mi.(Model)
	if m.cursor != 1 {
		t.Fatalf("after one down, cursor = %d, want 1", m.cursor)
	}

	mi, _ = m.Update(keyMsg("down"))
	m = mi.(Model)
	mi, _ = m.Update(keyMsg("down"))
	m = mi.(Model)
	if m.cursor != 2 {
		t.Fatalf("cursor should clamp at last index, got %d, want 2", m.cursor)
	}

	mi, _ = m.Update(keyMsg("up"))
	m = mi.(Model)
	mi, _ = m.Update(keyMsg("up"))
	m = mi.(Model)
	mi, _ = m.Update(keyMsg("up"))
	m = mi.(Model)
	if m.cursor != 0 {
		t.Fatalf("cursor should clamp at zero, got %d, want 0", m.cursor)
	}
}

func TestListStateTransitions(t *testing.T) {
	withSelection := newTestModel([]storage.Challenge{{ID: "1"}})
	empty := newTestModel(nil)

	cases := []struct {
		name      string
		m         Model
		key       string
		wantState sessionState
	}{
		{"a opens add form", withSelection, "a", stateAdd},
		{"e opens edit form with selection", withSelection, "e", stateEdit},
		{"e does nothing without selection", empty, "e", stateList},
		{"d opens delete confirm with selection", withSelection, "d", stateDeleteConfirm},
		{"d does nothing without selection", empty, "d", stateList},
		{"enter opens detail view with selection", withSelection, "enter", stateDetail},
		{"v opens detail view with selection", withSelection, "v", stateDetail},
		{"v does nothing without selection", empty, "v", stateList},
		{"t opens ctf filter prompt", withSelection, "t", stateFilterCTFPrompt},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mi, _ := tc.m.Update(keyMsg(tc.key))
			got := mi.(Model)
			if got.state != tc.wantState {
				t.Errorf("after %q, state = %v, want %v", tc.key, got.state, tc.wantState)
			}
		})
	}
}

func TestListFilterToggles(t *testing.T) {
	challenges := []storage.Challenge{
		{ID: "1", Status: storage.StatusUnsolved, FlaggedHard: false},
		{ID: "2", Status: storage.StatusMissed, FlaggedHard: true},
	}

	t.Run("f toggles flagged-hard filter", func(t *testing.T) {
		m := newTestModel(challenges)
		mi, _ := m.Update(keyMsg("f"))
		m = mi.(Model)
		if !m.filters.FlaggedHard {
			t.Fatal("expected FlaggedHard filter to be enabled")
		}
		if len(m.filtered) != 1 || m.filtered[0].ID != "2" {
			t.Fatalf("filtered = %v, want only challenge 2", m.filtered)
		}
	})

	t.Run("m toggles missed-only filter", func(t *testing.T) {
		m := newTestModel(challenges)
		mi, _ := m.Update(keyMsg("m"))
		m = mi.(Model)
		if !m.filters.MissedOnly {
			t.Fatal("expected MissedOnly filter to be enabled")
		}
		if len(m.filtered) != 1 || m.filtered[0].ID != "2" {
			t.Fatalf("filtered = %v, want only challenge 2", m.filtered)
		}
	})

	t.Run("c cycles status filter through all states", func(t *testing.T) {
		m := newTestModel(challenges)
		want := []string{
			string(storage.StatusUnsolved),
			string(storage.StatusSolved),
			string(storage.StatusMissed),
			"",
		}
		for i, w := range want {
			mi, _ := m.Update(keyMsg("c"))
			m = mi.(Model)
			if m.filters.Status != w {
				t.Fatalf("cycle %d: status filter = %q, want %q", i, m.filters.Status, w)
			}
		}
	})

	t.Run("x clears all filters", func(t *testing.T) {
		m := newTestModel(challenges)
		m.filters = FilterSettings{CTF: "foo", Status: string(storage.StatusMissed), FlaggedHard: true, MissedOnly: true}
		mi, _ := m.Update(keyMsg("x"))
		m = mi.(Model)
		if m.filters != (FilterSettings{}) {
			t.Fatalf("filters = %+v, want zero value", m.filters)
		}
		if len(m.filtered) != len(challenges) {
			t.Fatalf("filtered = %v, want all challenges restored", m.filtered)
		}
	})
}

func TestListInstantToggles(t *testing.T) {
	t.Run("h toggles flagged hard without touching status", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1", Status: storage.StatusUnsolved, FlaggedHard: false}})
		mi, cmd := m.Update(keyMsg("h"))
		m = mi.(Model)
		if !m.challenges[0].FlaggedHard {
			t.Fatal("expected FlaggedHard to be toggled on")
		}
		if m.challenges[0].Status != storage.StatusUnsolved {
			t.Fatalf("status should be untouched, got %v", m.challenges[0].Status)
		}
		if cmd == nil {
			t.Fatal("expected a save command to be returned")
		}
	})

	t.Run("s cycles status without touching flagged hard", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1", Status: storage.StatusUnsolved, FlaggedHard: false}})
		mi, cmd := m.Update(keyMsg("s"))
		m = mi.(Model)
		if m.challenges[0].Status != storage.StatusSolved {
			t.Fatalf("status = %v, want Solved", m.challenges[0].Status)
		}
		if m.challenges[0].FlaggedHard {
			t.Fatal("FlaggedHard should be untouched by status cycling")
		}
		if cmd == nil {
			t.Fatal("expected a save command to be returned")
		}
	})
}

func TestDeleteConfirm(t *testing.T) {
	t.Run("y confirms and returns to list with a delete command", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1"}})
		m.state = stateDeleteConfirm
		mi, cmd := m.Update(keyMsg("y"))
		m = mi.(Model)
		if m.state != stateList {
			t.Fatalf("state = %v, want stateList", m.state)
		}
		if cmd == nil {
			t.Fatal("expected a delete command to be returned")
		}
	})

	t.Run("n cancels back to list without a command", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1"}})
		m.state = stateDeleteConfirm
		mi, cmd := m.Update(keyMsg("n"))
		m = mi.(Model)
		if m.state != stateList {
			t.Fatalf("state = %v, want stateList", m.state)
		}
		if cmd != nil {
			t.Fatal("expected no command when cancelling")
		}
	})

	t.Run("esc cancels back to list", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1"}})
		m.state = stateDeleteConfirm
		mi, _ := m.Update(keyMsg("esc"))
		m = mi.(Model)
		if m.state != stateList {
			t.Fatalf("state = %v, want stateList", m.state)
		}
	})
}

func TestDetailStateKeys(t *testing.T) {
	t.Run("esc/q/v/enter return to list", func(t *testing.T) {
		for _, key := range []string{"esc", "q", "v", "enter"} {
			m := newTestModel([]storage.Challenge{{ID: "1"}})
			m.state = stateDetail
			mi, _ := m.Update(keyMsg(key))
			m = mi.(Model)
			if m.state != stateList {
				t.Errorf("key %q: state = %v, want stateList", key, m.state)
			}
		}
	})

	t.Run("e opens edit form", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1"}})
		m.state = stateDetail
		mi, _ := m.Update(keyMsg("e"))
		m = mi.(Model)
		if m.state != stateEdit {
			t.Fatalf("state = %v, want stateEdit", m.state)
		}
	})

	t.Run("h toggles flagged hard only", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1", Status: storage.StatusUnsolved, FlaggedHard: false}})
		m.state = stateDetail
		mi, cmd := m.Update(keyMsg("h"))
		m = mi.(Model)
		if !m.challenges[0].FlaggedHard {
			t.Fatal("expected FlaggedHard to be toggled on")
		}
		if m.challenges[0].Status != storage.StatusUnsolved {
			t.Fatalf("status should be untouched, got %v", m.challenges[0].Status)
		}
		if cmd == nil {
			t.Fatal("expected a save command to be returned")
		}
	})

	t.Run("s cycles status without touching flagged hard", func(t *testing.T) {
		m := newTestModel([]storage.Challenge{{ID: "1", Status: storage.StatusUnsolved, FlaggedHard: false}})
		m.state = stateDetail
		mi, cmd := m.Update(keyMsg("s"))
		m = mi.(Model)
		if m.challenges[0].Status != storage.StatusSolved {
			t.Fatalf("status = %v, want Solved", m.challenges[0].Status)
		}
		if m.challenges[0].FlaggedHard {
			t.Fatal("s should not toggle FlaggedHard as a side effect")
		}
		if cmd == nil {
			t.Fatal("expected a save command to be returned")
		}
	})
}

func TestMessageHandlers(t *testing.T) {
	t.Run("challengesLoadedMsg populates challenges and filtered", func(t *testing.T) {
		m := Model{state: stateList}
		challs := []storage.Challenge{{ID: "1"}, {ID: "2"}}
		mi, _ := m.Update(challengesLoadedMsg{challenges: challs})
		m = mi.(Model)
		if len(m.challenges) != 2 || len(m.filtered) != 2 {
			t.Fatalf("challenges/filtered not populated: %+v", m)
		}
	})

	t.Run("errMsg sets errorMsg and clears infoMsg", func(t *testing.T) {
		m := Model{state: stateList, infoMsg: "old info"}
		mi, _ := m.Update(errMsg{err: errBoom})
		m = mi.(Model)
		if m.errorMsg != errBoom.Error() {
			t.Fatalf("errorMsg = %q, want %q", m.errorMsg, errBoom.Error())
		}
		if m.infoMsg != "" {
			t.Fatalf("infoMsg should be cleared, got %q", m.infoMsg)
		}
	})

	t.Run("operationSuccessMsg sets infoMsg and returns to list", func(t *testing.T) {
		m := Model{state: stateAdd, errorMsg: "old error"}
		mi, cmd := m.Update(operationSuccessMsg("done"))
		m = mi.(Model)
		if m.infoMsg != "done" {
			t.Fatalf("infoMsg = %q, want %q", m.infoMsg, "done")
		}
		if m.errorMsg != "" {
			t.Fatalf("errorMsg should be cleared, got %q", m.errorMsg)
		}
		if m.state != stateList {
			t.Fatalf("state = %v, want stateList", m.state)
		}
		if cmd == nil {
			t.Fatal("expected a reload command to be returned")
		}
	})

	t.Run("cancelFormMsg returns to list", func(t *testing.T) {
		m := Model{state: stateAdd}
		mi, _ := m.Update(cancelFormMsg{})
		m = mi.(Model)
		if m.state != stateList {
			t.Fatalf("state = %v, want stateList", m.state)
		}
	})

	t.Run("applyCTFFilterMsg sets filter and returns to list", func(t *testing.T) {
		m := Model{state: stateFilterCTFPrompt, challenges: []storage.Challenge{{ID: "1", CTFName: "HTB"}, {ID: "2", CTFName: "picoCTF"}}}
		mi, _ := m.Update(applyCTFFilterMsg{ctf: "HTB"})
		m = mi.(Model)
		if m.filters.CTF != "HTB" {
			t.Fatalf("filters.CTF = %q, want %q", m.filters.CTF, "HTB")
		}
		if m.state != stateList {
			t.Fatalf("state = %v, want stateList", m.state)
		}
		if len(m.filtered) != 1 || m.filtered[0].ID != "1" {
			t.Fatalf("filtered = %v, want only challenge 1", m.filtered)
		}
	})
}

type staticErr string

func (e staticErr) Error() string { return string(e) }

const errBoom = staticErr("boom")

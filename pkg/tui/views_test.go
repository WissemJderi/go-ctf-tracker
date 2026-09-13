package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/WissemJderi/go-ctf-tracker/pkg/storage"
	"github.com/charmbracelet/lipgloss"
)

func TestListControlsPanelSticksToBottom(t *testing.T) {
	challs := []storage.Challenge{
		{ID: "1", CTFName: "HTB", Name: "Alpha", Points: 10, Difficulty: "easy", Status: storage.StatusUnsolved},
		{ID: "2", CTFName: "HTB", Name: "Beta", Points: 20, Difficulty: "medium", Status: storage.StatusSolved},
	}
	height := 24
	m := Model{state: stateList, height: height, width: 100, filtered: challs}

	view := m.View()
	lines := strings.Split(view, "\n")

	if n := len(lines); n != height {
		t.Fatalf("View() rendered %d lines for height %d, want %d\nfull view:\n%s", n, height, height, view)
	}

	controlsH := lipgloss.Height(strings.TrimRight(m.listHelpView(), "\n"))
	controlsStart := height - controlsH

	if !strings.Contains(lines[height-1], "quit tracker") {
		t.Fatalf("last line must be the bottom of the controls panel, got %q", lines[height-1])
	}
	found := false
	for i := controlsStart; i < height; i++ {
		if strings.Contains(lines[i], "CONTROLS") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("controls header must sit within the last %d lines, got %q", controlsH, lines[controlsStart])
	}
}

func TestListControlsPanelPinnedWithManyChallenges(t *testing.T) {
	challs := make([]storage.Challenge, 100)
	for i := range challs {
		challs[i] = storage.Challenge{
			ID:         fmt.Sprintf("%d", i),
			CTFName:    "HTB",
			Name:       fmt.Sprintf("challenge-%02d", i),
			Points:     10,
			Difficulty: "easy",
			Status:     storage.StatusUnsolved,
		}
	}
	height := 24
	m := Model{state: stateList, height: height, width: 100, filtered: challs, cursor: 95}

	view := m.View()
	lines := strings.Split(view, "\n")

	if n := len(lines); n != height {
		t.Fatalf("View() rendered %d lines for height %d with 100 challenges, want %d\nfull view:\n%s", n, height, height, view)
	}
	if !strings.Contains(lines[height-1], "quit tracker") {
		t.Fatalf("controls must stay pinned to the bottom with many challenges, got %q", lines[height-1])
	}
}

func TestListViewScrollsWindowToCursor(t *testing.T) {
	challs := make([]storage.Challenge, 100)
	for i := range challs {
		challs[i] = storage.Challenge{
			ID:         fmt.Sprintf("%d", i),
			CTFName:    "HTB",
			Name:       fmt.Sprintf("challenge-%02d", i),
			Points:     10,
			Difficulty: "easy",
			Status:     storage.StatusUnsolved,
		}
	}

	t.Run("cursor near top shows first rows", func(t *testing.T) {
		m := Model{state: stateList, height: 24, width: 100, filtered: challs, cursor: 1}
		view := m.View()
		if !strings.Contains(view, "challenge-01") {
			t.Fatalf("expected visible challenge-01 at top of list, full view:\n%s", view)
		}
		if strings.Contains(view, "challenge-99") {
			t.Fatalf("did not expect visible challenge-99 near top of list, full view:\n%s", view)
		}
	})

	t.Run("cursor near bottom shows last rows", func(t *testing.T) {
		m := Model{state: stateList, height: 24, width: 100, filtered: challs, cursor: 98}
		view := m.View()
		if !strings.Contains(view, "challenge-98") {
			t.Fatalf("expected visible challenge-98 near bottom of list, full view:\n%s", view)
		}
		if strings.Contains(view, "challenge-0") {
			t.Fatalf("did not expect visible challenge-0 near bottom of list, full view:\n%s", view)
		}
	})
}

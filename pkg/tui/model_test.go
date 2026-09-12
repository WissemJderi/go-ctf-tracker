package tui

import (
	"testing"

	"github.com/WissemJderi/go-ctf-tracker/pkg/storage"
)

func TestRotateStatus(t *testing.T) {
	cases := []struct {
		name string
		in   storage.ChallengeStatus
		want storage.ChallengeStatus
	}{
		{"unsolved to solved", storage.StatusUnsolved, storage.StatusSolved},
		{"solved to missed", storage.StatusSolved, storage.StatusMissed},
		{"missed to unsolved", storage.StatusMissed, storage.StatusUnsolved},
		{"unknown to unsolved", storage.ChallengeStatus("bogus"), storage.StatusUnsolved},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rotateStatus(tc.in)
			if got != tc.want {
				t.Errorf("rotateStatus(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNormalizeCursor(t *testing.T) {
	three := make([]storage.Challenge, 3)
	empty := []storage.Challenge{}

	cases := []struct {
		name     string
		filtered []storage.Challenge
		cursor   int
		want     int
	}{
		{"empty list resets to zero", empty, 5, 0},
		{"cursor within bounds unchanged", three, 1, 1},
		{"cursor past end clamps to last index", three, 10, 2},
		{"negative cursor clamps to zero", three, -1, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeCursor(nil, tc.filtered, tc.cursor)
			if got != tc.want {
				t.Errorf("normalizeCursor(_, %d items, cursor=%d) = %d, want %d", len(tc.filtered), tc.cursor, got, tc.want)
			}
		})
	}
}

func TestApplyFilters(t *testing.T) {
	challenges := []storage.Challenge{
		{ID: "1", CTFName: "HTB Cyber Apocalypse", Status: storage.StatusUnsolved, FlaggedHard: false},
		{ID: "2", CTFName: "HackTheBox Weekly", Status: storage.StatusSolved, FlaggedHard: true},
		{ID: "3", CTFName: "picoCTF", Status: storage.StatusMissed, FlaggedHard: true},
		{ID: "4", CTFName: "picoCTF", Status: storage.StatusUnsolved, FlaggedHard: false},
	}

	cases := []struct {
		name    string
		filters FilterSettings
		wantIDs []string
	}{
		{
			name:    "no filters returns everything",
			filters: FilterSettings{},
			wantIDs: []string{"1", "2", "3", "4"},
		},
		{
			name:    "CTF name filter is case-insensitive substring match",
			filters: FilterSettings{CTF: "hackthebox"},
			wantIDs: []string{"2"},
		},
		{
			name:    "status filter matches exact status",
			filters: FilterSettings{Status: string(storage.StatusMissed)},
			wantIDs: []string{"3"},
		},
		{
			name:    "flagged hard filter keeps only flagged challenges",
			filters: FilterSettings{FlaggedHard: true},
			wantIDs: []string{"2", "3"},
		},
		{
			name:    "missed only filter keeps only missed challenges",
			filters: FilterSettings{MissedOnly: true},
			wantIDs: []string{"3"},
		},
		{
			name:    "filters combine with AND semantics",
			filters: FilterSettings{CTF: "picoCTF", FlaggedHard: true},
			wantIDs: []string{"3"},
		},
		{
			name:    "no matches returns empty slice",
			filters: FilterSettings{CTF: "does-not-exist"},
			wantIDs: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := Model{challenges: challenges, filters: tc.filters}
			got := m.applyFilters()

			gotIDs := make([]string, 0, len(got))
			for _, c := range got {
				gotIDs = append(gotIDs, c.ID)
			}

			if len(gotIDs) != len(tc.wantIDs) {
				t.Fatalf("applyFilters() = %v, want %v", gotIDs, tc.wantIDs)
			}
			for i := range gotIDs {
				if gotIDs[i] != tc.wantIDs[i] {
					t.Fatalf("applyFilters() = %v, want %v", gotIDs, tc.wantIDs)
				}
			}
		})
	}
}

package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ChallengeStatus string

const (
	StatusUnsolved ChallengeStatus = "Unsolved"
	StatusSolved   ChallengeStatus = "Solved"
	StatusMissed   ChallengeStatus = "Missed"
)

type Challenge struct {
	ID          string          `json:"id"`
	CTFName     string          `json:"ctf_name"`
	Name        string          `json:"name"`
	Category    string          `json:"category"`
	Points      int             `json:"points"`
	Difficulty  string          `json:"difficulty"` // e.g., Easy, Medium, Hard
	Status      ChallengeStatus `json:"status"`     // Unsolved, Solved, Missed
	FlaggedHard bool            `json:"flagged_hard"`
	Notes       string          `json:"notes"`
	WriteupURL  string          `json:"writeup_url"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type JSONStorage struct {
	filepath string
}

func NewJSONStorage(customPath string) (*JSONStorage, error) {
	var path string
	if customPath != "" {
		path = customPath
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			path = "db.json" // Fallback to current directory
		} else {
			dir := filepath.Join(home, ".config", "ctf-tracker")
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}
			path = filepath.Join(dir, "db.json")
		}
	}

	return &JSONStorage{filepath: path}, nil
}

func (s *JSONStorage) GetPath() string {
	return s.filepath
}

func (s *JSONStorage) Load() ([]Challenge, error) {
	if _, err := os.Stat(s.filepath); os.IsNotExist(err) {
		return []Challenge{}, nil
	}

	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %w", err)
	}

	if len(data) == 0 {
		return []Challenge{}, nil
	}

	var challenges []Challenge
	if err := json.Unmarshal(data, &challenges); err != nil {
		return nil, fmt.Errorf("failed to parse database JSON: %w", err)
	}

	// Validate loaded data
	for i, ch := range challenges {
		if ch.ID == "" {
			return nil, fmt.Errorf("challenge at index %d has empty ID", i)
		}
		if ch.CTFName == "" {
			return nil, fmt.Errorf("challenge at index %d has empty CTF name", i)
		}
		if ch.Name == "" {
			return nil, fmt.Errorf("challenge at index %d has empty name", i)
		}
		if ch.Status == "" {
			return nil, fmt.Errorf("challenge at index %d has empty status", i)
		}
		switch ch.Status {
		case StatusUnsolved, StatusSolved, StatusMissed:
		default:
			return nil, fmt.Errorf("challenge at index %d has invalid status: %s", i, ch.Status)
		}
	}

	return challenges, nil
}

func (s *JSONStorage) Save(challenges []Challenge) error {
	data, err := json.MarshalIndent(challenges, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	// Use atomic write to avoid race conditions and orphaned temp files
	tempFile := s.filepath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tempFile, s.filepath); err != nil {
		// Clean up temp file on failure
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

func (s *JSONStorage) Add(ch Challenge) error {
	challenges, err := s.Load()
	if err != nil {
		return err
	}

	if ch.ID == "" {
		ch.ID = generateShortID()
	}
	ch.CreatedAt = time.Now()
	ch.UpdatedAt = time.Now()

	challenges = append(challenges, ch)
	return s.Save(challenges)
}

func (s *JSONStorage) Update(ch Challenge) error {
	challenges, err := s.Load()
	if err != nil {
		return err
	}

	found := false
	for i, c := range challenges {
		if c.ID == ch.ID {
			ch.UpdatedAt = time.Now()
			challenges[i] = ch
			found = true
			break
		}
	}

	if !found {
		return errors.New("challenge not found")
	}

	if err := s.Save(challenges); err != nil {
		return err
	}
	return nil
}

func (s *JSONStorage) Delete(id string) error {
	challenges, err := s.Load()
	if err != nil {
		return err
	}

	index := -1
	for i, c := range challenges {
		if c.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("challenge not found")
	}

	challenges = append(challenges[:index], challenges[index+1:]...)
	if err := s.Save(challenges); err != nil {
		return err
	}
	return nil
}

func (s *JSONStorage) Get(id string) (Challenge, error) {
	challenges, err := s.Load()
	if err != nil {
		return Challenge{}, err
	}

	for _, c := range challenges {
		if c.ID == id {
			return c, nil
		}
	}

	return Challenge{}, errors.New("challenge not found")
}

func generateShortID() string {
	for attempt := 0; attempt < 10; attempt++ {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			// Fallback to time-based unique ID in the extremely rare case of rand failure
			return fmt.Sprintf("%x", time.Now().UnixNano())[:8]
		}
		id := hex.EncodeToString(b)
		if id == "" {
			continue
		}
		return id
	}
	// Last resort: use timestamp-based ID with nanoseconds
	return fmt.Sprintf("%x", time.Now().UnixNano())[:8]
}

package storage

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestJSONStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	// 1. Initial load should return empty slice
	challs, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(challs) != 0 {
		t.Errorf("Expected 0 challenges, got %d", len(challs))
	}

	// 2. Add challenge
	ch := Challenge{
		CTFName:    "Test CTF",
		Name:       "Super Crypto",
		Category:   "Crypto",
		Points:     100,
		Difficulty: "Easy",
		Status:     StatusUnsolved,
	}

	err = storage.Add(ch)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// 3. Load again and verify
	challs, err = storage.Load()
	if err != nil {
		t.Fatalf("Load after Add failed: %v", err)
	}
	if len(challs) != 1 {
		t.Fatalf("Expected 1 challenge, got %d", len(challs))
	}

	savedCh := challs[0]
	if savedCh.Name != "Super Crypto" || savedCh.Points != 100 || savedCh.ID == "" {
		t.Errorf("Saved challenge fields mismatch: %+v", savedCh)
	}

	// 4. Update challenge
	savedCh.Status = StatusSolved
	savedCh.FlaggedHard = true
	savedCh.WriteupURL = "https://writeup.local"

	err = storage.Update(savedCh)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updatedCh, err := storage.Get(savedCh.ID)
	if err != nil {
		t.Fatalf("Get after Update failed: %v", err)
	}
	if updatedCh.Status != StatusSolved || !updatedCh.FlaggedHard || updatedCh.WriteupURL != "https://writeup.local" {
		t.Errorf("Updated fields mismatch: %+v", updatedCh)
	}

	// 5. Delete challenge
	err = storage.Delete(savedCh.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	challs, err = storage.Load()
	if err != nil {
		t.Fatalf("Load after Delete failed: %v", err)
	}
	if len(challs) != 0 {
		t.Errorf("Expected 0 challenges after deletion, got %d", len(challs))
	}
}

func TestJSONStorage_Validation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	// Test empty file handling
	if err := storage.Save([]Challenge{}); err != nil {
		t.Errorf("Save of empty slice failed: %v", err)
	}

	// Test loading from empty file
	challs, err := storage.Load()
	if err != nil {
		t.Fatalf("Load from empty file failed: %v", err)
	}
	if len(challs) != 0 {
		t.Errorf("Expected empty slice from empty file, got %d", len(challs))
	}

	// Test invalid JSON
	invalidJSON := `{"id": "", "ctf_name": "", "name": ""}`
	if err := os.WriteFile(dbPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	_, err = storage.Load()
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	// Test invalid status
	invalidStatusJSON := `[{"id": "test", "ctf_name": "Test", "name": "Test Challenge", "category": "Misc", "points": 100, "difficulty": "Easy", "status": "InvalidStatus", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}]`
	if err := os.WriteFile(dbPath, []byte(invalidStatusJSON), 0644); err != nil {
		t.Fatalf("Failed to write invalid status JSON: %v", err)
	}

	_, err = storage.Load()
	if err == nil {
		t.Error("Expected error for invalid status, got nil")
	}

	// Test missing required fields
	missingFieldsJSON := `[{"id": "test", "ctf_name": "Test", "category": "Misc", "points": 100, "difficulty": "Easy"}]`
	if err := os.WriteFile(dbPath, []byte(missingFieldsJSON), 0644); err != nil {
		t.Fatalf("Failed to write missing fields JSON: %v", err)
	}

	_, err = storage.Load()
	if err == nil {
		t.Error("Expected error for missing fields, got nil")
	}
}

func TestJSONStorage_UpdateNonExistent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	ch := Challenge{
		CTFName:    "Test CTF",
		Name:       "Test Challenge",
		Category:   "Misc",
		Points:     100,
		Difficulty: "Easy",
		Status:     StatusUnsolved,
	}

	err = storage.Add(ch)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Try to update non-existent challenge
	err = storage.Update(Challenge{
		ID:   "nonexistent",
		Name: "NonExistent",
	})
	if err == nil {
		t.Error("Expected error updating non-existent challenge, got nil")
	}

	if err.Error() != "challenge not found" {
		t.Errorf("Expected 'challenge not found' error, got: %v", err)
	}
}

func TestJSONStorage_DeleteNonExistent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	// Try to delete non-existent challenge
	err = storage.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error deleting non-existent challenge, got nil")
	}

	if err.Error() != "challenge not found" {
		t.Errorf("Expected 'challenge not found' error, got: %v", err)
	}
}

func TestJSONStorage_ConcurrentAccess(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	// Add initial challenge
	if err := storage.Add(Challenge{
		CTFName:    "Test CTF",
		Name:       "Test Challenge",
		Category:   "Misc",
		Points:     100,
		Difficulty: "Easy",
		Status:     StatusUnsolved,
	}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Test concurrent reads
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := storage.Load()
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestGenerateShortID(t *testing.T) {
	generated := make(map[string]bool)
	maxAttempts := 1000

	for i := 0; i < maxAttempts; i++ {
		id := generateShortID()
		if id == "" {
			t.Error("generateShortID returned empty string")
		}
		if len(id) != 8 {
			t.Errorf("generateShortID returned string of length %d, expected 8", len(id))
		}
		if generated[id] {
			t.Logf("ID collision detected on attempt %d", i)
		}
		generated[id] = true
	}

	// With 4 bytes (32 bits), we expect about 1 collision in 10k attempts
	// In 1000 attempts, collisions are rare, which is expected
	t.Logf("Generated %d unique IDs out of %d attempts", len(generated), maxAttempts)
}

func TestNewJSONStorage_DefaultPathUnderHome(t *testing.T) {
	fakeHome, err := os.MkdirTemp("", "ctf-tracker-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(fakeHome)

	t.Setenv("HOME", fakeHome)

	storage, err := NewJSONStorage("")
	if err != nil {
		t.Fatalf("NewJSONStorage failed: %v", err)
	}

	wantPath := filepath.Join(fakeHome, ".config", "ctf-tracker", "db.json")
	if storage.GetPath() != wantPath {
		t.Errorf("GetPath() = %q, want %q", storage.GetPath(), wantPath)
	}

	if info, err := os.Stat(filepath.Join(fakeHome, ".config", "ctf-tracker")); err != nil || !info.IsDir() {
		t.Errorf("expected config directory to be created, stat error: %v", err)
	}
}

func TestJSONStorage_LoadMissingFileReturnsEmptySlice(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "does-not-exist.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	challs, err := storage.Load()
	if err != nil {
		t.Fatalf("Load on missing file returned error: %v", err)
	}
	if challs == nil || len(challs) != 0 {
		t.Errorf("expected empty non-nil slice, got %v", challs)
	}
}

func TestJSONStorage_LoadEmptyJSONArray(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	if err := os.WriteFile(dbPath, []byte("[]"), 0644); err != nil {
		t.Fatalf("failed to write empty array: %v", err)
	}

	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	challs, err := storage.Load()
	if err != nil {
		t.Fatalf("Load of empty JSON array failed: %v", err)
	}
	if len(challs) != 0 {
		t.Errorf("expected 0 challenges, got %d", len(challs))
	}
}

func TestJSONStorage_SaveLeavesNoTempFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	if err := storage.Save([]Challenge{{ID: "1", CTFName: "CTF", Name: "Chal", Status: StatusUnsolved}}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(dbPath + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("expected temp file to be removed after successful save, stat error: %v", err)
	}
}

func TestJSONStorage_SaveFailureLeavesNoTempFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Point the db at a path whose parent directory doesn't exist, so the
	// rename step fails and Save must clean up its temp file.
	dbPath := filepath.Join(tempDir, "missing-parent", "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	err = storage.Save([]Challenge{{ID: "1", CTFName: "CTF", Name: "Chal", Status: StatusUnsolved}})
	if err == nil {
		t.Fatal("expected Save to fail when parent directory is missing")
	}

	if _, err := os.Stat(dbPath + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("expected temp file to be cleaned up after failed save, stat error: %v", err)
	}
}

func TestJSONStorage_AddPreservesCallerSuppliedID(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	if err := storage.Add(Challenge{ID: "custom-id", CTFName: "CTF", Name: "Chal", Status: StatusUnsolved}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	ch, err := storage.Get("custom-id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if ch.ID != "custom-id" {
		t.Errorf("ID = %q, want %q", ch.ID, "custom-id")
	}
}

func TestJSONStorage_AddSetsTimestamps(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	before := time.Now()
	if err := storage.Add(Challenge{CTFName: "CTF", Name: "Chal", Status: StatusUnsolved}); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	after := time.Now()

	challs, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	ch := challs[0]

	if ch.CreatedAt.Before(before) || ch.CreatedAt.After(after) {
		t.Errorf("CreatedAt = %v, want between %v and %v", ch.CreatedAt, before, after)
	}
	if diff := ch.UpdatedAt.Sub(ch.CreatedAt); diff < 0 || diff > time.Second {
		t.Errorf("CreatedAt and UpdatedAt should both be set on creation, close together: %v vs %v", ch.CreatedAt, ch.UpdatedAt)
	}
}

func TestJSONStorage_UpdateOnlyTouchesTargetChallenge(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	if err := storage.Add(Challenge{ID: "a", CTFName: "CTF", Name: "First", Status: StatusUnsolved}); err != nil {
		t.Fatalf("Add(a) failed: %v", err)
	}
	if err := storage.Add(Challenge{ID: "b", CTFName: "CTF", Name: "Second", Status: StatusUnsolved}); err != nil {
		t.Fatalf("Add(b) failed: %v", err)
	}

	a, err := storage.Get("a")
	if err != nil {
		t.Fatalf("Get(a) failed: %v", err)
	}
	a.Status = StatusSolved
	if err := storage.Update(a); err != nil {
		t.Fatalf("Update(a) failed: %v", err)
	}

	b, err := storage.Get("b")
	if err != nil {
		t.Fatalf("Get(b) failed: %v", err)
	}
	if b.Name != "Second" || b.Status != StatusUnsolved {
		t.Errorf("challenge b was unexpectedly modified: %+v", b)
	}
}

func TestJSONStorage_DeleteOnlyRemovesTargetChallenge(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	for _, id := range []string{"a", "b", "c"} {
		if err := storage.Add(Challenge{ID: id, CTFName: "CTF", Name: id, Status: StatusUnsolved}); err != nil {
			t.Fatalf("Add(%s) failed: %v", id, err)
		}
	}

	if err := storage.Delete("b"); err != nil {
		t.Fatalf("Delete(b) failed: %v", err)
	}

	challs, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(challs) != 2 {
		t.Fatalf("expected 2 remaining challenges, got %d", len(challs))
	}
	for _, c := range challs {
		if c.ID == "b" {
			t.Fatal("challenge b should have been removed")
		}
	}
}

func TestJSONStorage_GetOnEmptyStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ctf-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test-db.json")
	storage, err := NewJSONStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	_, err = storage.Get("anything")
	if err == nil {
		t.Fatal("expected error getting from empty storage, got nil")
	}
	if err.Error() != "challenge not found" {
		t.Errorf("expected 'challenge not found' error, got: %v", err)
	}
}

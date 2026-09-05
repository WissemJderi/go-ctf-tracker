package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WissemJderi/go-ctf-tracker/pkg/storage"
	"github.com/WissemJderi/go-ctf-tracker/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dbPath := flag.String("db", "", "Path to the JSON database file (defaults to ~/.config/ctf-tracker/db.json)")
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Println("CTF Challenge Tracker v1.0.0")
		return
	}

	// Initialize Storage
	store, err := storage.NewJSONStorage(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Validate storage path
	if _, err := os.Stat(store.GetPath()); os.IsNotExist(err) {
		// Allow new files to be created, but warn if path is not writable
		dir := filepath.Dir(store.GetPath())
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating database directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize Model
	model := tui.NewModel(store)

	// Run Program in AltScreen mode for seamless fullscreen user experience
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI program: %v\n", err)
		os.Exit(1)
	}
}

# ⚡ Go CTF Tracker (ctf-tracker) ⚡

A lightning-fast Terminal UI (TUI) application written in Go to track Capture The Flag (CTF) challenges with minimal fuss. Built for live events and post-CTF review — fast, portable, and keyboard-driven.

---

## 🚀 Key Features

- ⚡ Blazing fast & lightweight: written in pure Go with no heavy database dependencies. Data is stored locally in a single JSON file (`~/.config/ctf-tracker/db.json`).
- 🔥 Flag Hard challenges: mark difficult or time-consuming tasks during a live CTF so you can revisit them later.
- 🔗 Writeup links & notes: save writeup URLs and notes per challenge, then mark them resolved after reviewing the writeup.
- 📊 Status & difficulty tracking: manage challenges across statuses (`Unsolved`, `Solved`, `Missed`) and difficulties (`Easy`, `Medium`, `Hard`).
- 🔍 Instant filtering: filter by CTF name, status, flagged hard (`🔥`), or show missed-only challenges (`✖`).
- ⌨️ Interactive TUI: built with Charmbracelet (`bubbletea`, `lipgloss`) for a responsive, keyboard-first terminal experience.

---

## 📦 Installation

Prerequisites: Go 1.22+.

```bash
# Clone the repository
git clone https://github.com/WissemJderi/go-ctf-tracker.git
cd go-ctf-tracker

# Build the binary
go build -o ctf-tracker .

# (Optional) Move to your PATH
sudo mv ctf-tracker /usr/local/bin/
```

---

## 🎮 Usage

Run the TUI:

```bash
./ctf-tracker
```

Use a custom database file:

```bash
./ctf-tracker -db ./my-ctfs.json
```

---

## ⌨️ Keyboard Controls

| Key / Shortcut | Action |
| :--- | :--- |
| `a` | Add a new challenge |
| `e` | Edit the selected challenge |
| `d` | Delete the selected challenge |
| `v` or `Enter` | View full challenge details (Notes & Writeup URL) |
| `s` | Cycle challenge status (`Unsolved` ➔ `Solved` ➔ `Missed` ➔ `Unsolved`) |
| `h` | Toggle **Flagged Hard** (`🔥`) |
| `t` | Filter by CTF name (prompt) |
| `c` | Cycle status filter (`All` ➔ `Unsolved` ➔ `Solved` ➔ `Missed`) |
| `f` | Toggle **Hard only** filter (`🔥`) |
| `m` | Toggle **Missed only** filter (`✖`) |
| `x` | Clear all filters |
| `q` or `Ctrl+C` | Quit tracker |

---

## 💡 Suggested Live CTF Workflow

1. During a live CTF:
   - Press `a` to quickly add challenges as you find them.
   - Mark blockers or time sinks with `h` so teammates can skip and return later.
2. After the CTF / when writeups are available:
   - Use `f` (Hard only) or `m` (Missed only) to focus on challenges to follow up.
   - Select a challenge, press `v` or `Enter` to view details, then `e` to paste writeup URLs and add notes.
   - Press `s` to mark a challenge as `Solved` after reviewing the writeup.

This keeps your live decisions recorded and turns them into actionable post-game follow-ups.

---

## 🤝 Contributing

Contributions, suggestions, and bug reports are welcome!

1. Fork the repository.
2. Create a feature branch: `git checkout -b my-feature`.
3. Make your changes and verify with `go build`.
4. Open a pull request describing your changes.

For larger changes, please open an issue first to discuss the design.

---

## 🛡️ License

This project is licensed under the MIT License — see the `LICENSE` file for details.

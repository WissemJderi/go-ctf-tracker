# ⚡ Go CTF Tracker (`ctf-tracker`) ⚡

A blazing-fast Terminal UI (TUI) application written in **Go** to track CTF (Capture The Flag) challenges effortlessly. Designed for speed, simplicity, and efficiency—whether you are playing solo or with a team.

---

## 🚀 Key Features

- **⚡ Blazing Fast & Lightweight**: Written in pure Go with zero heavy database dependencies; stores your challenges locally in a clean, portable JSON format (`~/.config/ctf-tracker/db.json`).
- **🔥 Live Hard Challenge Flagging**: Flag difficult or time-consuming challenges during a live CTF so you can easily find them later.
- **🔗 Writeup Linking & Post-CTF Resolution**: When writeups are published after a CTF, revisit your flagged or missed challenges, paste the writeup URL, take notes, and mark them as solved!
- **📊 Comprehensive Status Management**: Track challenges across multiple statuses (`Unsolved`, `Solved`, `Missed`) and difficulties (`Easy`, `Medium`, `Hard`).
- **🔍 Advanced Filtering**: Filter your view instantly by CTF Name, Status, Flagged Hard (`🔥`), or Missed-Only challenges (`✖`).
- **⌨️ Fully Interactive TUI**: Built with Charmbracelet (`bubbletea`, `lipgloss`) for a gorgeous, responsive, keyboard-driven terminal experience.

---

## 📦 Installation

Make sure you have Go 1.22+ installed on your system.

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

Run the TUI by simply executing:

```bash
./ctf-tracker
```

You can also specify a custom database file path:
```bash
./ctf-tracker -db ./my-ctfs.json
```

---

## ⌨️ Keyboard Controls & Keybindings

| Key / Shortcut | Action |
| :--- | :--- |
| `a` | Add a new CTF challenge |
| `e` | Edit the selected challenge |
| `d` | Delete the selected challenge |
| `v` or `Enter` | View full challenge details (Notes & Writeup URL) |
| `s` | Cycle challenge status (`Unsolved` ➔ `Solved` ➔ `Missed` ➔ `Unsolved`) |
| `h` | Toggle **Flagged Hard** (`🔥`) instantly |
| `t` | Filter by CTF Name prompt |
| `c` | Cycle Status filter (`All` ➔ `Unsolved` ➔ `Solved` ➔ `Missed`) |
| `f` | Toggle **Hard Only** filter (`🔥`) |
| `m` | Toggle **Missed Only** filter (`✖`) |
| `x` | Clear all active filters |
| `q` or `Ctrl+C` | Quit tracker |

---

## 💡 Live CTF Workflow Example

1. **During the Live CTF**: 
   - Press `a` to quickly add challenges as you discover them.
   - If you encounter a brutal crypto or reversing challenge that stumps your team, press `h` (or check `Flagged Hard?`) to mark it with a 🔥.
2. **After the CTF / Writeups Released**:
   - Filter your view by pressing `f` (Hard Only) or `m` (Missed Only).
   - Press `v` or `Enter` to open the details view, press `e` to edit, paste the author's GitHub/blog writeup into the **Writeup URL** field, add key takeaways to **Notes**, and press `s` to mark it as `Solved`!

---

## 🛡️ License

MIT License. Feel free to fork, contribute, and open source!

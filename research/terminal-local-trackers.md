# What other terminal and local-first CTF trackers offer

Research for issue #2 of go-ctf-tracker. Question: what do existing terminal and
local-first CTF trackers offer, and what floor of expectations does that set for
a solo player? This file lists the directly adjacent tools, the adjacent
workflows people actually run in terminals, the features that repeat across
tools, the features that appear once, and the takeaway for go-ctf-tracker.

Method: web and GitHub searches for CTF-specific CLI/TUI trackers, then a wider
net over task and note CLIs plus git and markdown workflows that show up in CTF
contexts. Only tools that a single player could run locally are covered in
depth. Organizer-side authoring tools get a short pass for contrast.

## The short answer

Very few players face a dedicated offline tracker, because almost nobody
published one. The repeatable pattern across what exists is small and
consistent. Track each challenge, keep a solved flag, keep a notes field, link
or generate the writeup, show progress per CTF and category, and keep it all in
one local file you own. go-ctf-tracker already covers status, notes, writeup
URL, filters, and the single JSON file. The clearly missing pieces are storing
the flag value itself, a solved timestamp, hint usage, a per-CTF workspace, and
a progress summary.

## Direct peers: terminal and local-first CTF trackers

### CTFMATE
URL: https://github.com/Sycosmile/CTFMATE
Language: Python CLI.
Feature set: auto-created folder tree per challenge, flag tracker, hints-used
count, notes per challenge, live dashboard, writeup template generator, multi-CTF
with switch command, overall stats with category solve rates.
Philosophy: the terminal is the workspace. It generates the folder structure so
a player never loses recon, exploits, loot, or tools, and it turns tracking into
a few commands.
What a solo player gets: a full event workspace from `new` to `writeup`, with
flags, hints, and notes captured along the way.
Deliberately skips: no TUI, no offline scoreboard, no timer (on the roadmap), no
search (on the roadmap).

### htbConsole
URL: https://pypi.org/project/htbconsole
Language: Python, Textual TUI plus a headless CLI.
Feature set: browse and filter challenges by category and difficulty, start/stop
containers, submit flags, download and auto-extract task files, lazy-loaded
community writeups, per-challenge notes, per-CTF work for the HTB CTF platform
with live scoreboard, and CLI one-shot queries that print JSON or YAML.
Philosophy: wrap the whole HTB and CTFd-style platform in a keyboard UI so the
player never leaves the terminal.
What a solo player gets: the reference for what a CTF TUI layout looks like.
Filter, detail panel, submit flag, notes per task.
Deliberately skips: offline use entirely. It is a thin client over the HTB API.
Note: the closest existing example of a real CTF TUI, and its workflow (list,
filter, detail, notes, submit) is almost exactly go-ctf-tracker's shape.

### CTFdCLI
URL: https://github.com/Yeeb1/CTFdCLI
Language: Python CLI.
Feature set: sync challenges from a CTFd instance, download and organize files,
generate a README per challenge, submit flags from a folder or the command line,
track solved status and attempts, multi-profile for several competitions,
offline-first after the sync.
Philosophy: pull the platform down to the filesystem, then work locally.
What a solo player gets: a local mirror of the event with progress tracking and
flag submission from the terminal.
Deliberately skips: no TUI, no leaderboard, no stats beyond solved counts.

### ctf_thief_tools
URL: https://github.com/Kochanac/ctf_thief_tools
Language: small Python and bash scripts.
Feature set: dump and update challenges from a handful of platforms (CTFd,
ctf-forces, innoctf, and others) into JSON, list tasks, and mark one solved with
`setattr Solved True`.
Philosophy: the minimum viable tracker. Scripts that talk to the platform API
and flip a status field.
What a solo player gets: a lightweight solved-or-not view without entering the
browser.
Deliberately skips: notes, hints, writeups, categories beyond the dump, anything
except status.

### CTFdown
URL: https://github.com/sbeving/CTFdown
Language: Python CLI.
Feature set: download challenges from a CTFd platform, organize into folders by
category, write a Markdown file per challenge with name, value, category,
description, tags, and file links.
Philosophy: offline archiving first. Turn the live platform into a browsable,
readable, permanent local copy.
What a solo player gets: a local study archive for training after the event.
Deliberately skips: status tracking, notes, flags, writeups. It stops at the
archive.

### HackLog (0xAquila ctf-tracker)
URL: https://github.com/0xAquila/ctf-tracker
Language: single-file HTML app, vanilla JavaScript, localStorage. Not a
terminal app, included because it shares the local-first rail.
Feature set: per-challenge sessions with live timer, methodology checklists for
Recon/Exploitation/Privesc/Documentation plus per-category sublists, a findings
logger for flags, credentials, hashes, CVEs, and ports, cheat sheets, a writeup
journal, notes, XP, levels, activity heatmap, three themes, and JSON export and
import.
Philosophy: one file, zero dependencies, everything in your browser, export your
data any time. Marketing language leans hard on offline and ownership.
What a solo player gets: the most complete local tracker found, including stuff
go-ctf-tracker does not have, like flag capture, timestamps on findings, stats,
and per-phase checklists.
Deliberately skips: any network or sync, any dependency, any real backend.

Context on the authoring side, for contrast only. ctfcli (CTFd's own CLI),
nsec/ctf-script, and ctf-term are for people who build challenges, not solve
them. They manage challenge-as-code, deploy services, and run local training
engines with leaderboards and hint penalties. They show the same terminal-first
instinct but a different audience, and none of them track a solo player's solve
progress.

## Git and markdown workflows used in CTF contexts

Before and alongside these tools, players track with plain files. This is the
default there is no tool for, and go-ctf-tracker is competing with it.

Writeup repos. The dominant pattern is a git repo per CTF or per player with one
Markdown README per challenge and a script or two beside it (oranav/ctf-writeups,
evyatar9/Writeups, yrudwls/ctf-writeup, snwau/picoCTF writeups, and many more).
The structure is CTF, then category, then challenge. Progress is expressed as
which README files exist and what the git log shows. This pattern assumes the
writeup is the record.

Status-folder writeups. Some players organize by state instead of category,
keeping a done/ folder and an in-progress/ folder with a table of challenge,
status, and writeup path (RoryPoonyth/ctf-writeups). The status table is the
tracker; the folders are the evidence.

Loot-style logs. At least one player keeps a running REPORT.md with PENDING and
pre-solved sections and stores found flags as a named loot list, including flags
that were rejected or retired (0xsaju/ctf). This is flag capture as a habit,
which no general note tool enforces.

Obsidian vaults and note templates. Shared Obsidian vaults for CTF exist to
prescribe the same structure a tracker would: per-category challenge notes, a
progress dashboard, reusable templates for investigation steps, evidence,
flags, and lessons learned (dwightdoran/ctf-note-template). HTTPSB note tools go
further and auto-generate a machine note and vault structure from the platform
API (0x4xel/HTNotes). These vaults live on local Markdown files, so they are
local-first without being terminal-first.

## Adjacent terminal task and note tools

These are not CTF-specific but show up in the same terminal workflows, and the
task bar for a tracker borrows from them. With one exception (org-mode's
progress logging), none of them appeared in CTF-specific searches. Players reach
for them because they are already installed, not because they track CTFs.

Taskwarrior. C++ CLI task manager. Tasks carry metadata, arbitrary tags, due
dates, project, and reports with filters. It logs nothing about a challenge but
gives a fast, keyboard-driven list, which is the competence floor. Taskwarrior
is a close workflow match for "one key to flip state, filter, report".

todo.txt. One task per line in a plain text file. Its whole thesis is a file you
can read, edit, and diff without the tool. That thesis matches go-ctf-tracker's
single JSON file: the data format is the product.

Org-mode. Emacs outline with TODO states and progress logging. When you move a
task to DONE, org can record the timestamp and an optional note (org-todo with a
prefix, or automatic CLOSED timestamps). This is the one adjacent tool that
routinely appears in CTF contexts, because org files show up inside writeup and
note workflows. Timestamped state changes are exactly what a post-event timeline
review wants, and org treats them as a first-class feature.

orgwarrior. A Go CLI that manages org-mode files as a task list. add, done,
modify, delete, list, filters by tag and date, relative dates. Proof that Go +
org-file backends are a live combination in the local tracker space.

Joplin terminal. Markdown note CLI with notebooks, tags, search, and sync to
several backends. Vim-style and command-line modes. A general contributor to the
"notes in the terminal" habit, not a CTF-specific one.

nb. Plaintext note CLI backed by plain folders and git. Notebooks, tags, search.
Again general, but its line "the plain text note is the archive" is the same
philosophy as the writeup-repo pattern.

vimwiki. Personal wiki inside Vim. Same local Markdown wiki idea as Obsidian but
editable in the terminal.

## Features that appear repeatedly

Repetition is weak evidence of need. These appeared across the largest set of
tools:

- Solved or attempt status per challenge, toggled with one command. CTFMATE
  solve, ctf_thief_tools setattr, CTFdCLI status columns, htbConsole submit,
  HackLog sessions. go-ctf-tracker has this.
- A per-challenge notes field. CTFMATE, htbConsole, HackLog, the Obsidian
  templates, go-ctf-tracker. The most repeated feature of all.
- Storing the flag value itself. CTFMate flag tracker, HackLog findings logger,
  CTFdCLI and htbConsole submissions, and the REPORT.md loot pattern. Players
  treat the flag as the artifact to keep.
- A writeup slot. Writeup URL (go-ctf-tracker), writeup templates (CTFMATE),
  writeup journal (HackLog), generated per-challenge README (CTFdCLI), and the
  entire writeup-repo pattern.
- Category, difficulty, and points metadata. CTFMATE, CTFdCLI, htbConsole,
  ctf-term, CTFdown, HackLog. go-ctf-tracker has all three.
- Filtering by category, status, or difficulty. htbConsole, CTFdCLI, ctf-term,
  go-ctf-tracker.
- Multi-event management. CTFMATE new/switch, CTFdCLI multi-profile, htbConsole
  CTF list, HackLog sessions, CTFdown folders per CTF. go-ctf-tracker models
  this only as a ctf_name field on every entry plus a filter prompt.
- A progress overview. CTFMATE dashboard and stats, CTFdCLI "1 solved (4.0%)",
  htbConsole profile stats, HackLog dashboard and heatmap, ctf-term stats, the
  Obsidian progress dashboard. go-ctf-tracker has no counts or summary.
- Timestamps on state changes. org-mode progress logging, HackLog timestamped
  findings, go-ctf-tracker created/updated, and git history in the writeup-repo
  pattern. Missing in go-ctf-tracker for the solve event specifically.
- Hint usage tracking. CTFMATE hint count, ctf-term hints with penalty, HackLog
  hint tooltips. Present but thinner than the rest.
- Local, portable, human-readable data. Single JSON (go-ctf-tracker), JSON
  export (HackLog), plain Markdown (repositories, Obsidian, CTFdown).
  go-ctf-tracker's single JSON file already fits this floor.
- Platform sync. CTFdCLI, htbConsole, ctf_thief_tools, CTFdown all pull from
  the platform API. This repeats only inside the connected-tools cluster, never
  in the offline-solo cluster.

## Features that appear once

One-offs, listed so nobody mistakes them for a floor.

- AI assistant with full session context. HackLog.
- XP, levels, achievements, activity heatmap. HackLog.
- Phase-based methodology checklists. HackLog.
- Live per-challenge timer. HackLog (CTFMATE lists it only on its roadmap).
- Automatic workspace folder scaffolding. CTFMATE, and on the HTB side HTNotes.
- Connection info extraction from challenge text (host, port, URL). CTFdCLI.
- Auto flag submission from flag.txt sitting in the challenge folder. CTFdCLI.
- Referencing challenges by index number instead of name. CTFMATE.
- First blood bonuses, hint penalties, and scoreboards. ctf-term, but that is
  an organizer-side engine.
- Full-text search across challenges. htbConsole and ctf-term. Borderline: it
  appeared twice, and go-ctf-tracker only filters rather than searches.

## Floor of expectations

A solo player who keeps their CTF workflow in the terminal walks in with these
expectations, because the tools and file workflows above repeat them:

1. Add while the event runs, in one or two keystrokes.
2. Flip status instantly: solved, missed, working. One key, no form.
3. A notes field on every challenge.
4. A place to store the flag the moment they find one, because losing a flag is
   the failure state every tool exists to prevent.
5. A writeup slot linked to the challenge, either a URL or a generated scaffold.
6. A view of progress per CTF and per category: counts and simple stats, not a
   raw list.
7. Data in one local file that is readable and portable without the tool.
8. Timestamps that let a player rebuild what happened and when after the event.

## What go-ctf-tracker clearly lacks

Against that floor, the clear gaps, biggest first:

- Storing the flag value itself. Every peer tracks the flag; go-ctf-tracker only
  tracks a solved state and a notes string.
- A progress summary. No counts per status, per category, or per CTF, and no
  dashboard view.
- A per-CTF workspace. The only grouping is a text filter prompt over one flat
  list, so a player cannot isolate one event or switch context cleanly.
- A solved timestamp. Only created and updated are stored, so the timeline of
  solves cannot be reconstructed.
- Hint usage tracking. CTFMATE and ctf-term treat this as a core field.
- Connection information (host, port, URL). Present in the CTFd world through
  CTFdCLI and platform metadata, and needed for any challenge that runs a
  service.
- Full-text search. Filters exist but there is no search across challenge names
  and notes.
- Writeup scaffolding. The writeup URL field exists, which covers the link
  half; nothing generates a writeup template or a journal.
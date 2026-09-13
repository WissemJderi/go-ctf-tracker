# Web platforms, dashboards, and note templates: what they set as tracker expectations

Research for decision: which features go-ctf-tracker should add, and which it should leave to the big platforms.

Scope: CTFd (competition platform a player actually logs into), CTFtime (event calendar plus result/writeup archive), the practice platforms (HackTheBox, TryHackMe, picoCTF, OverTheWire), and the note-template ecosystem (Notion, Obsidian, spreadsheets, markdown trackers) that solo players bolt on top of all of them.

All sources were read only. This document records what each one tracks, the expectations each sets, and what each already solves so well that a lightweight local tool should not re-implement it.

---

## 1. CTFd

CTFd is the most common competition platform a solo player actually logs into during an event. The platform is the event.

What it tracks for players:

- A challenge list with name, category, point value, description, attached files, hints, and flags. Challenges can be locked or hidden until a threshold is met (unlockable challenges).
- Correct/incorrect flag submissions, including after a challenge is already solved. The platform tells you a submission is right or wrong.
- Category (web, crypto, forensics, pwn, reverse, misc) per challenge. A few events add tags on top.
- Per-challenge solve state. Solved challenges drop out or grey out; the scoreboard updates in real time with tie resolution.
- Dynamic scoring for many events: points decay as more teams solve a challenge (initial, decay, minimum, function columns).
- First-blood indicators: the scoreboard can show first capture, and organizers see time-to-first-capture from the capture feed. The public scoreboard generally does not surface per-challenge first blood unless the theme adds it.
- Solutions visibility: newer CTFd lets organizers store a solution per challenge and reveal visibility after the event.
- Challenge ratings and reviews: upvotes/downvotes after solving, plus admin-only reviews.
- Own past submissions, if the organizer enables that view.

What it tracks for organizers:

- User and team management: roles, bans, hidden teams, invitations, team avatars.
- The scoreboard with freeze and hide scores at set times, score graphs comparing the top teams, and a progression matrix (per player or per team, per challenge) on the statistics page.
- A capture log feed for statistics such as time to first capture.
- Export and import of the whole CTF (challenges, teams, configs) as a zip. CSV exports of the scoreboard, including per-challenge splits.
- Tracking events: first time a player opens a challenge, hint unlocks.

Expectations it sets for players:

- "Solved" is a binary online state, computed by the platform from flag submission. The tracker does not need to decide solved; the platform is the source of truth during an event.
- Category and points are platform-owned metadata. A player should not have to hand-type them; they are one fetch away.
- Hints, files, and remote connection info live on the challenge card and follow the player around during the event.

What it solves so thoroughly that go-ctf-tracker should NOT duplicate:

- Scoring, dynamic scoring, leaderboards, ties, freeze. A local tracker has no league sense and should not fake one.
- Flag validation. The platform is authoritative; an offline tracker cannot know if a flag string is correct.
- Submission history and first-blood ordering. It is the event's own record, and it goes away with the platform at the end of the event.
- Challenge vault hosting the hints, files, and descriptions while the event runs.

The platform-native hole: when the event ends, the platform closes or is torn down (static export is a manual, optional step organizers often skip). Everything you did is claimable only from the scoreboard and your own notes. Post-event, CTFd gives you nothing. That is the gap.

---

## 2. CTFtime

CTFtime is the aggregation and archive layer. Players browse it before an event and return to it after.

What it offers players:

- Upcoming events: name, format (Jeopardy, Attack-Defense, mixed), on-line or on-site location, UTC start and end, duration, weight, and a "will participate" tick so teams and players broadcast which events they plan to play.
- A filterable event list (recent, upcoming, past, high-school, academic, finals, by weight) with a feed and an iCalendar export so events land in your calendar.
- Team ratings and points, season by season, with a formula and weight voting. Global, country, and year rankings.
- Scores and full standings for past events, including per-team points and rating points gained.
- Per-task solve presence: the scoreboard feed accepts per-task points and solve timestamps, so standings can show which tasks a team solved and when.
- Writeup aggregation: a task page per challenge, and a team page listing the writeups its players posted, each with a link. There is a global writeups list and team writeup counts.
- Per-team and per-player pages with participation history: what events they played, how many points, and which tasks they wrote up. Team comparison page and a teams directory.
- An API exporting results, teams, and events in JSON, extensible by third-party tools (several MCP servers and CLIs wrap it).

Expectations it sets for players:

- An event calendar and result history should live outside any single event platform. Registration, scheduling, and retrospective results are event-independent.
- Ratings and standings are a team-level public concern, kept sane by a weight formula and voting, not a personal tracking chore.
- Writeups are expected after the event and are stored for a long time as a shared resource. CTFtime's own FAQ says most team websites die and writeups vanish, so it asks players to post a copy there.

What it solves so thoroughly that go-ctf-tracker should NOT duplicate:

- Upcoming event discovery, scheduling, and the ratings ecosystem. All team- and community-scale; useless to a solo offline tool.
- Global standings, team ratings, weights, and the writeup index that is public by design.

Its holes: CTFtime knows a task and a resolution roughly, but not your attempt, your notes, or your flag. It aggregates writeups; it does not host working notes, failed paths, or a private journal. And it only knows events registered with it. Everything in the middle, from "event is on my calendar" to "writeup posted", is the player's own work.

---

## 3. Practice platforms

Practice platforms behave like permanent gymnasiums. Their competitive surface (leaderboards, points) matters less to a solo player than their per-challenge progress tracking, which is exactly the model most trackers imitate.

### 3.1 HackTheBox

What it tracks:

- Machines and challenges with OS (Linux/Windows/other), difficulty (Easy to Insane), release date, machine rating, user rating, and count of user/system owns.
- Per-machine flags split by user and root: `user.txt` and `root.txt`, both MD5, owned independently. Machines show "user own" and "system/root own" counts separately.
- Owned state, in-progress state, a favorites list, a todo list, and a dashboard sliced into favorites, in progress, and recommended items.
- Points, skill badges, achievements, activity feed, and profile stats including content completed.
- Official writeups (VIP) and a community writeup ecosystem cross-referenced by technique, OS, difficulty, CVE, and certification relevance.
- Timed flags: seasons, fortresses (many flags on one host), challenges (single flag), and Sherlock DFIR investigations with multiple flags or tasks.
- Rankings by users, teams, countries, and universities, plus per-machine annual leagues.

Expectations it sets for players:

- A challenge object has structured metadata: name, OS, difficulty, points, category, and its own official writeup link. Players copy this into their trackers.
- Progress is partial by design. A machine can be "user owned" but not "root owned"; a fortress has flags 3 of 11. Solved is not binary along a single axis; it is a list of captured sub-objectives.
- A queue mechanism (todo list, favorites, recommended) exists to decide what to attempt next, separate from what is already owned.
- Re-visiting is a first-class idea: retired machines stay worth doing, and official writeups are meant for review after an attempt.

What it solves so thoroughly that go-ctf-tracker should NOT duplicate:

- The metadata database for every machine and challenge (names, OS, difficulty, flags paths, writeup links). Keeping a copy of this by hand is waste; any tool should treat HTB's data as read-only reference.
- Official writeup hosting and community indexes. Publishing a redone walkthrough is redundant.
- Owned-state bookkeeping on the platform itself (todo, favorites, owned splits). It works while you are logged in and is the source of truth for the game.

Its holes: the account marks owned, but not what you learned. No place for your note trail, your creds, your failed paths, or why a technique worked. Writeups are someone else's; your working notes are not stored. And the valuable metadata (retired machines list, per-machine technique tags) is behind login or VIP.

### 3.2 TryHackMe

What it tracks:

- Rooms grouped into learning paths and modules, each with difficulty color (green beginner, yellow intermediate, red advanced/insane), a room type (walkthrough versus challenge/CTF), and per-room tag categories.
- Per-question progress inside a room, plus completion percentage for paths and modules.
- Points earned per room, with separate rules for walkthrough and challenge rooms and for monthly versus all-time leaderboards.
- A streak counter for consecutive days of engagement, with streak freezes from weekly missions.
- A capability score derived from rooms completed across difficulty tiers plus a consistency/decay component, mapped to level names.
- First-completion bonuses (the first player through a room gets extra points), monthly and all-time leaderboards.
- For team and business users: assignments with due dates, time reports, activity reports, a skills matrix built from room tags, and CSV export of task progress.

Expectations it sets for players:

- Difficulty is a coarse, color-coded tier, not a precise number, and it drives both points and perceived progression.
- Content is sketched as a path, not a flat free-for-all. Players expect to see "what is next" ordered by topic and difficulty.
- Long-term consistency is tracked (streaks, decay, capability score) and surfaced as motivation; engagement is itself a measured outcome.
- A room can be partially done: mid-way through a 20-question room, the platform records exactly where you stopped, so you can resume.

What it solves so thoroughly that go-ctf-tracker should NOT duplicate:

- Detailed per-question progress across hundreds of rooms, path enrollment, badges, and points economy. That is THM's entire product.
- Streaks and monthly leaderboard grind. Artificial engagement loops are the opposite of an offline personal tool.

Its holes: the platform tracks completion, not learning. Nothing captures what you struggled with, what the room taught you, or a comparison across platforms ("I do THM for intro web, HTB for AD"). Per-question progress lives only inside THM and evaporates if a room retires (old walkthrough rooms stop counting toward your score entirely).

### 3.3 picoCTF and picoGym

What it tracks:

- Challenges listed by points (50 to 500 historically, roughly tracking difficulty), split into fixed categories: General Skills, Cryptography, Forensics, Web Exploitation, Reverse Engineering, Binary Exploitation, plus a newer Blockchain category.
- Competition events per year, plus picoGym, a permanent practice space holding challenges from past competitions that players can solve year-round.
- Team score combined from challenge solves, awarded once per challenge. Individual solve still counts for personal practice.
- Classrooms for teachers with per-student stats per event, category progress graphs compared to the class average, a separate classroom scoreboard, and CSV export.
- Playlists and learning resources (picoPrimer, video tutorials) that group challenges around a topic.
- Solves and a score counter that ticks up in the corner when you submit a correct flag.

Expectations it sets for players:

- Competitions and practice are the same challenge pool; the same problems come back to solve at your own pace. Archiving an event's challenges should be standard behavior.
- Category is the primary mental model, and learning paths are organized around it, sequenced easy to hard.
- Points equal difficulty well enough to triage: low point values are read as approachable.

What it solves so thoroughly that go-ctf-tracker should NOT duplicate:

- The free, permanent challenge archive with per-category, per-point filtering. Re-listing picoGym's catalog inside a local tool adds nothing.
- Classroom analytics and the education ecosystem.

Its holes: nothing records your notes, your attempts, or your per-category growth over time. The gym tells you what you solved; it cannot tell you what you are actually good at or what you should restudy. The 2026 picoCTF guide notes that tags disagree with the real technique more often than expected, so a learned-then-tagged local record beats the platform's own labels after the fact.

### 3.4 OverTheWire (Bandit and the wargames)

What it tracks:

- The platform itself tracks nothing. Bandit gives one numbered page per level with a goal and sometimes a hint; the level progression is linear and the passwords form the progress.
- OverTheWire states plainly: passwords are not saved automatically, start over if you lose them, and take notes because levels get harder and you will want to return or help others.

Expectations it sets for players:

- Progress is a single scalar per game: current level, rank, and number of levels solved. Bandit is structured as an ordered ladder, unlike HTB's free-for-all.
- Because the platform keeps no state at all, every player ends up building their own password vault plus per-level notes. Many public CLIs exist exactly for this (otw-cli saves passwords and per-level markdown; BanditCLI auto-detects passwords and marks a level complete).

What it avoids that go-ctf-tracker should copy:

- OverTheWire has no dashboard, no tags, no difficulty ratings, no score graph. It is the purest evidence that a linear progression tracker of highest-level-per-game plus one note per level is a complete, expected pattern. This is the smallest viable tracker shape, and players already treat it as normal.
- Its hole is the same as everything above and is the platform's point: the record is entirely yours, and laborious without a tool.

---

## 4. Note-template ecosystems

This is where solo players hand-build the tracking layer the platforms do not provide. The templates are the externalized answer to "what fields matter".

### 4.1 The recurring fields

Collecting the templates in circulation (Notion marketplace CTF trackers and notes templates, the widely shared Elite Era writeup template, CTF-Citadel challenge tables, HeroCTF/hackyeaster writeup repos, the ctf-writeup skill's submission format, Kioptrix/Obsidian lab-note guides):

- Challenge or machine name. Always present.
- Status. Solved, in progress, or not started; the writeup-repo tables add "ready/100%, done, solved, relaunched, research" granular states. One repo uses a status flow fetched -> seen -> working -> hoard -> solved.
- Category or tags. Web, pwn, crypto, forensics, reverse, misc, plus per-platform tags (OSINT, AD, cloud).
- Difficulty. Either a word (easy/medium/hard/insane) or a number/color; writes up every time.
- Points. Present whenever the platform publishes them (picoCTF, CTFd dynamic scores).
- Flag. Stored, generally redacted in shared repos, with a flag format note.
- Writeup URL or link. In nearly every table, either the official writeup or the player's own writeup file.
- Date. Solve date, event dates, or release date of the box.
- Notes. Free text; the richer templates split it into goal, key clues, plan, steps, solution summary, and lessons learned.
- Platform or event name. picoCTF, HackTheBox, TryHackMe, the specific CTF. Cross-platform trackers always keep it.
- Time spent. Common but optional; used for triage and writeup effort estimation.

Sub-objectives appear only where the platform has them: HTB templates store user flag and root flag and OS in separate fields, and the Obsidian machine-tracker plugin pulls 40+ fields from the HTB API (name, OS, difficulty plus playtime and resolves) into a per-machine YAML note.

### 4.2 Template workflows

- Obsidian players keep one folder per challenge or machine, one markdown note per phase (recon, enumeration, exploitation, proof, lessons learned), frontmatter carrying the metadata, and wiki-links between notes. The guidance is unanimous on this shape: headings and metadata are fixed and reused for every box, notes must be replayable and timestamped, and failed paths are kept in compressed form because they are diagnostic memory, not clutter.
- Spreadsheet players (Sheets, Excel) keep rows of challenges and columns of the recurring fields. Little else; the spreadsheet is the tracker's minimum viable pivot table.
- Notion players run a database with the same columns plus a status select and a rollup of writeups. The marketplace templates demonstrate that people want sorting and filtering by category, difficulty, and OS (one HTB template review: "sorting every machine with its tags, difficulties and os"), plus a template button that stamps a detail page per challenge.
- Writeup repos on GitHub follow the same schema as tracker rows: a per-challenge markdown file plus a README table of challenge, category, difficulty, points, status, and a link to the writeup. The ctf-writeup agent skill enforces a frontmatter of title, ctf, date, category, difficulty, points, flag format, author, then a key-observation, single-script, flag document.

### 4.3 What the template ecosystem expects

- A tracker row is reusable metadata: name, category, difficulty, points, status, flag, date, writeup link, platform, notes. Those are the fields nearly every template repeats; a tool that stores all of them for one challenge satisfies the ecosystem's vocabulary.
- Writeups are a separate, later artifact built from the same metadata, and turning live notes into a writeup should be a light step (the writeup builders and note guides all promise "reports are mostly a rephrasing of good notes, not a new creative project").
- The ecosystem already made the portability call: markdown and JSON win, clouds lose. Obsidian, leaflet, m-ctf, CTFx, HackLog, and several others run on plain files in a local folder.

### 4.4 What the template ecosystem solves that go-ctf-tracker should NOT duplicate

- The template ecosystem is not a feature target; fighting Obsidian's graph or Notion's database view is lost effort and contradicts the local, keyboard-driven, single-file premiss. Its contribution is the field vocabulary and the writeup-from-notes flow, which are the right model to match, and the pain evidence below, which is the right problem to solve.

---

## 5. The expectations solo players carry in

Summing across all four source families:

1. A challenge is a typed row with name, category, tags, difficulty, points, flag, and a writeup link. Every platform and every template agrees; players expect to see exactly these fields, sorted and filtered.
2. Solved state is visible at a glance and is decided by the platform while the event runs. The tracker's job is to record the outcome, plus notes, without re-inventing scoring or flag validation.
3. Progress is granular per objective. User and root flags owned separately (HTB), per-question rooms (THM), multi-flag fortresses, per-level passwords (Bandit): consumers expect to track partial progress on one challenge, not just done versus not done.
4. Difficulty drives triage and learning order. Every platform rates it, and templates copy it; players use it to pick what to attempt next.
5. A queue of what to try next is separate from the record of what is done. HTB's todo and favorites, THM's paths and weekly missions, picoGym playlists all build this in; players expect "waiting" and "recommended" as their own state, not a solved column.
6. Re-visiting and post-event review is expected to work. Retired machines, picoGym archives, and CTFtime writeups all assume you come back after the event; the writeup is the natural product of that revisit.
7. Data should be yours, portable, and alive after the event ends. The platforms expire, the accounts close, the rooms retire; markdown, JSON, and a single local file are what players converged on to make the record theirs.

---

## 6. The gaps the platforms leave

These are the concrete pain points a lightweight local tool has room to attack:

1. Nothing holds working notes during the solve. Platforms store outcomes; a live attempt's log of commands, creds, ports, hypotheses, and failed paths has no home. The strongest documentation pain evidence comes from the writeup-builder author: hours of documentation for every real solve, twenty to thirty percent of total CTF time.
2. The record is shattered across platforms and events. HTB knows your HTB, picoGym knows picoGym, CTFtime knows results, your spreadsheet knows nothing, and no single place answers "which web challenges have I done across every platform and what did I learn".
3. Everything dies or closes. CTFd instances go down, THM rooms retire and stop counting, team writeup sites vanish (CTFtime itself warns about this), and OverTheWire flat out does not save. Longevity is the recurring complaint, not feature starvation.
4. No per-category or per-topic learning history. Platforms track completion and points; none answer "what am I good at, what did I struggle with, what should I review". The picoCTF writeup guides argue platform tags are unreliable and the real technique only shows up after a re-tag. A local tool can hold the learned-technique tag a player assigns after the fact.
5. Offline play and offline records. Nothing in this ecosystem is usable without network, and several islands (HTB's metadata, some rooms) are paywalled or login-walled. The tracker is the one component that must work during a blacked-out flight or a sealed exam room.
6. Writeup generation is manual and slow. The templates and builders agree the final writeup should be assembled from structured notes, and players burn hours doing that assembly by hand. A tool that models the fields once and exports the writeup files for free removes the single most disliked part of the workflow.

---

## 7. Read for go-ctf-tracker

Expectations to meet: model one challenge as name, category, tags, difficulty, points, flag, status, date, writeup link, plus a notes block; allow partial progress per challenge; keep a queue separate from the completion record; make re-visit and post-event review natural; store everything in the one local JSON file and never require the network.

Spaces already owned by someone else, do not rebuild: scoring and live leaderboards (CTFd), the event calendar, ratings, and the public writeup index (CTFtime), the per-platform metadata and owned-state databases (HTB, THM, picoGym), and the general-purpose note-with-graph UX (Obsidian and friends).

The niche that opens, and where the tool should aim: the single local record that spans platforms and events, written while solving offline, structured with the recurring template fields, able to re-tag challenges with what was actually learned, able to resume a paused challenge from its own half-solved state, and able to emit writeup drafts from the same fields nobody else stores.
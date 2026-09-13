# Community pain points in CTF challenge tracking

Research ticket: "What do CTF players say they want and hate in tracking?" (issue #4)
Branch: research/community-pain-points
Date: 2026-09-13

## Method

I searched public web sources with websearch and webfetch: GitHub repositories where solo players publish their CTF note and writeup workflows, team student wikis, and one blog on writing CTF walkthroughs. Reddit and Discord are walled gardens here. Reddit returned only branding shells when fetched directly, so I used search indexes and the many public GitHub and blog mirrors of the same community. I recorded only content I could actually fetch and quote. I labeled each finding REPEATED (it appears in multiple independent sources) or one-off (single source).

## Findings

### 1. Tracking is scattered: scattered notes, random text files, or nothing at all
Tool makers who build for solo players describe the same starting point. The HackLog tracker README says most players "track their work in scattered notes, random text files, or nothing at all." The CTFMate README says "most CTF players keep challenges in messy folders." A whole genre of personal tools and Obsidian vault templates exists precisely to replace that scatter. This is a large number of players pre-tool, and it is the founding complaint behind every tracker I found.

- Source: 0xAquila/ctf-tracker README and Sycosmile/CTFMATE README
- Label: REPEATED

### 2. Flags and progress get lost between sessions
The same CTFMate README says players "lose flags between sessions." A solo player's writeup archive README (Ophirky/ctf-write-ups) lists "make it easier to resume challenges after a long break" as a reason for keeping the archive. Another solo player (RoryPoonyth/ctf-writeups) organizes the archive strictly by status (done/ vs in-progress/) and mentions "recovered local session history," which implies session continuity was broken. CTFNote exists so a team can permanently archive what was only living in ephemeral Discord channels.

- Source: Sycosmile/CTFMATE README, Ophirky/ctf-write-ups README, RoryPoonyth/ctf-writeups README, TFNS/CTFNote README
- Label: REPEATED

### 3. Writeups pile up and never get written
This is the most consistent complaint. CTFMate: players "never get around to writing proper writeups." The Ophirky archive lists dozens of picoCTF solves as "Most are undocumented," and marks several challenges "unfinished" or "unsolved" in the record. LearnHacking.io's walkthrough guide ends with "DON'T put off your write-up and never share!" and names "post-CTF fatigue, poor note-taking, or just writing fatigue" as the causes of incomplete writeups. The rctf organizer guide tells event runners to set writeup deadlines one to two weeks after the event and to offer incentives, which organizers only do because writeups do not arrive otherwise. Another player repo literally appends "(More writeups will be added as I continue solving challenges!)" to its index, a running backlog made visible.

- Source: Sycosmile/CTFMATE README, Ophirky/ctf-write-ups README, learnhacking.io walkthrough guide, rctf.osec.io "After the CTF" guide, t4mpr/ctf-writeups README
- Label: REPEATED

### 4. Players write writeups for their own learning and future reference, not just for publishing
The HackUCF team README states the purpose of their writeup repo is partly "we learn more by writing down what we have done and explaining it in depth." Ophirky's archive was started to "stay organized and track my learning journey" and "reflect on both successes and mistakes." LearnHacking.io says one reason to write is "to strengthen your own understanding," and the author reports "I frequently refer back to my own notes." Many solo repos describe themselves as personal tracks of progress. The writeup is a personal artifact first.

- Source: HackUCF/ctfs README, Ophirky/ctf-write-ups README, learnhacking.io walkthrough guide, multiple solo writeup repos
- Label: REPEATED

### 5. Challenge notes need a fixed structure: objective, what was tried, evidence, result, lesson
Community guidance repeatedly converges on the same fields for a per-challenge note. The StudSec student wiki says a challenge notepad should answer: what systems and protocols are relevant, what the challenge primitive is, what you tried so far and with what results, and what interesting information you found. A popular Obsidian CTF note template (dwightdoran/ctf-note-template) provides sections for objective, flag format, hints, targets, ports, credentials, chronological investigation in a thought/action/result format, evidence, final flag, solution steps, and key lesson. LearnHacking.io's writeup outline is initial info, enumeration, rabbit holes, path to solution, actual solution. Multiple solo writeup repos (amyy45, Zeeshan01001, t4mpr) describe near-identical formats: challenge description, approach, steps, final flag, key learnings, references.

- Source: wiki.studsec.nl CTFNote guide, dwightdoran/ctf-note-template README, learnhacking.io guide, several solo writeup repos
- Label: REPEATED

### 6. Solved / in-progress / blocked status is the core view during an event
The tools all lead with a status board. CTFNote shows a "list of all (unsolved) challenges" as its working view. The Obsidian template has a dashboard that tracks "in progress, blocked, or solved." A CTFd issue from organizers shows they too keep "track of solved/unsolved" per challenge. Players mirror this in repos with done/ and in-progress/ folders. Status is the primary thing a live tracker has to show.

- Source: wiki.studsec.nl CTFNote guide, dwightdoran/ctf-note-template README, CTFd issue #321, RoryPoonyth/ctf-writeups README
- Label: REPEATED

### 7. Time and consistency tracking should be automatic, not manual
The CTFStreakTracker sells itself as tracking practice time "on autopilot," after noting the CLI exists for "automatically track" and its extension exists so you avoid "manually pressing start/stop." HackLog starts each challenge with a live timer. CTFMate's roadmap lists a per-challenge timer as a planned feature. Three independent tools arrive at the same shape: the player wants the time recorded without the overhead of clocking in and out.

- Source: sainivedhh/CTFStreakTracker README, 0xAquila/ctf-tracker README, Sycosmile/CTFMATE README
- Label: REPEATED

### 8. Skills, stats, and streaks are wanted, usually as motivation
CTFStreakTracker renders contribution heatmaps and streak badges to the player's GitHub profile. HackLog has XP, levels, achievements, a 90-day heatmap, and skill rings. CTFMate has a stats view with category solve rates. These are three independent tools all shipping progression surfaces.

- Source: sainivedhh/CTFStreakTracker README, 0xAquila/ctf-tracker README, Sycosmile/CTFMATE README
- Label: REPEATED

### 9. Live notes should become the writeup with as little rework as possible
The Obsidian template says the chronological note "makes the note useful both during the event and when converting it into a write-up later." CTFMate generates structured markdown writeup templates from the tracked challenge. LearnHacking.io says the author picks Markdown notes "as I often use Markdown for the eventual blog post." Writeup generation is the natural bridge between the live record and the post-event deliverable, and players expect it.

- Source: dwightdoran/ctf-note-template README, Sycosmile/CTFMATE README, learnhacking.io guide
- Label: REPEATED

### 10. Event infrastructure disappears, so evidence must be captured during the event
LearnHacking.io warns that challenge files and services "won't stay up" after an event and therefore the player must collect payloads, files, screenshots, and traffic during the event, not after. The rctf organizer guide instructs runners to keep challenge services alive for only a few days so players can finish writeups, and to tear everything down after. The only durable record the player keeps is their own. For a solo player this makes the live capture step the foundation of any later writeup.

- Source: learnhacking.io guide, rctf.osec.io "After the CTF" guide
- Label: REPEATED

### 11. Flag capture is a first-class action, not a line item in a text file
CTFMate has a dedicated `solve` command that takes a challenge reference and the flag and marks it solved. HackLog's "Findings Logger" is introduced as "Never Lose a Flag Again" and timestamps flags and credentials as they are logged. The writeup guide says one goal of the writeup is to prove the flag was found. In a tracker, recording the flag is the event's closing transaction.

- Source: Sycosmile/CTFMATE README, 0xAquila/ctf-tracker README, learnhacking.io guide
- Label: REPEATED

### 12. Logging must cost as few keys as possible
The tools consistently minimize input overhead: HackLog promises "start a session in seconds" and a one-shot findings logger, CTFMate is built around single-line commands, CTFStreakTracker removes the start/stop step entirely, and CTFNote lets the whole team keep notes where they already chat (Discord). No source says this directly, but every tool fights the same enemy: the keystrokes spent recording are keystrokes not spent solving. I mark this as inferred from tool design rather than an explicit player quote.

- Source: 0xAquila/ctf-tracker, Sycosmile/CTFMATE, sainivedhh/CTFStreakTracker, TFNS/CTFNote (inferred from design)
- Label: one-off (design inference)

### 13. Hints used should be tracked
CTFMate tracks how many hints were used per challenge with a `hint` command. No other found source mentions hint tracking.

- Source: Sycosmile/CTFMATE README
- Label: one-off

## Workflows

### During the live event
The solo player opens a board showing every challenge in the event with solved/unsolved status (CTFNote, Templates). They pick a challenge and open its note, which starts with objective, flag format, hints, and any known targets, ports, or credentials (Templates, StudSec). As they work, they append short chronological entries in a thought/action/result shape: what they tried, what happened, what to try next (StudSec, Templates). Findings get logged the moment they appear: a flag, a credential, a hash, a port, a CVE, a screenshot (HackLog, CTFMate, LearnHacking). Some keep a scanner running in the background and note its results (purplestorm CTF-Notes). When solved, the player records the flag and the challenge flips to solved on the board (CTFMate, CTFNote). Time, if tracked at all, is expected to be automatic (CTFStreakTracker). The live note is deliberately skeletal: full prose comes later.

### Post-event writeup
The event infrastructure is usually already gone or about to be torn down, so the player works from the live notes, screenshots, and files they captured (LearnHacking, rctf guide). They convert each challenge note into a writeup, and because the note already holds the sequence of attempts, evidence, and the flag, the conversion is mostly ordering and light editing (Templates, CTFMate). Writeups land in a repo organized by event and category, or on a blog, or on CTFtime (SecurityInnovation, porters of the ctfs/community repo, LearnHacking). The writeup backlog is chronic; organizers set deadlines and prizes because spontaneous output is rare, and many players publish only a fraction of what they solved (rctf guide, Ophirky "Most are undocumented"). When the player later reviews, the archive exists so they can resume a challenge after a break or refind a technique (Ophirky, LearnHacking "refer back to my own notes").

## Feature asks

Map to what go-ctf-tracker could adopt:

1. Per-challenge status (solved, in progress, blocked, not started) as the primary board view
2. A structured per-challenge note with objective, flag format, hints, targets, credentials, and a chronological what-tried/what-result log
3. Fast, few-keystroke capture of findings (flag, credential, hash, port, CVE) while solving
4. Flag recording as a first-class action that also flips challenge status
5. Timestamped append-only log per challenge
6. Writeup generation from the live notes, exported to Markdown
7. Automatic session time tracking, or none at all (do not make the player clock in and out)
8. Hints-used counter per challenge
9. Stats per category and per event, plus optional streaks or heatmaps for motivation
10. Event-scoped grouping with an index that mirrors event/category organization
11. All data local and offline, with easy export/import (a single JSON file fits the observed solo preference)

## Gaps takeaway

The community converges on one arc: capture the challenge live with low overhead, keep a status board, and turn the captured record into a writeup later. The two chronic failures are the writeup backlog (post-event fatigue beats good intentions, so the record must be complete enough to finish cheaply later) and overhead during the event (every keystroke of logging is a keystroke not spent solving, so capture must be fast and automatic where possible). The strongest edge a small offline keyboard-driven tracker has is being fastest to the log line: no tab switching, no mouse, no clock-in ritual. The weakest parts commonly complained about in existing tools are exactly the heavyweight web dashboards and manual timers, so a TUI should avoid both.

## Evidence quality

Almost all findings are fetchable hard evidence: direct quotes from READMEs of personal tools and archives, a student wiki guide, a published walkthrough-writing blog, and an organizer's operational guide. These are strong on artifact and semi-strong on voice, they are players and teams describing their own recorded practice in public. The weak side is the absence of live Reddit and Discord thread content, which both blocked direct fetching, so I could not quote the crowd in the moment or count upvotes across threads. The one direct inference is the low-overhead principle (finding 12), which I flag as such. The blog, wiki, and READMEs all triangulate to the same conclusions, so I rate overall confidence moderate-to-high on the shape of the problem and its main pain points, and lower on which single features are most wanted, because the evidence base over-represents builders who already made a tracker and under-represents players who quietly quit tracking.
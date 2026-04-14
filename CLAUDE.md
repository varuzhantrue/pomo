# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build -o pomo .      # build the binary
go run . <command>      # run without building
go vet ./...            # lint
go test ./...           # run tests (none exist yet)
```

## Architecture

`pomo` is a zero-dependency Go CLI Pomodoro timer. All data is persisted locally.

**Entry point — [main.go](main.go)**  
Parses the subcommand (`start`, `break`, `long`, `status`) and delegates. `start` counts today's completed focus sessions first so the new session can be numbered correctly.

**Timer — [timer.go](timer.go)**  
`RunTimer(sessionType, totalSecs, sessionNum)` drives a tick loop with ANSI escape codes for an in-place progress bar. On SIGINT it prompts the user to save a partial session. On natural completion it calls `notifyDone` then `saveSession`.

**Storage — [storage.go](storage.go)**  
`Session` structs are appended to `~/.pomo/sessions.json` as a JSON array. `todaySessions()` filters to completed (non-partial) focus sessions starting after midnight local time.

**Notifications — [notify.go](notify.go)**  
Sends a terminal bell unconditionally; on macOS also fires a native notification via `osascript`. Errors are silently ignored (best-effort).

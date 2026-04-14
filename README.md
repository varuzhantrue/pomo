# pomo

A minimal Pomodoro timer for your terminal. No dependencies, no configuration, just focus.

## Features

- 25-minute focus sessions with auto-numbering
- 5-minute short breaks and 15-minute long breaks
- Live progress bar with color-coded countdown
- Saves session history locally — tracks total focus time per day
- Prompts to save partial sessions on Ctrl+C
- Native macOS notifications on session completion (terminal bell on all platforms)

## Install

**Prerequisites:** Go 1.21 or later.

```bash
git clone https://github.com/varuzhantrue/pomo.git
cd pomo
go build -o pomo .
```

Then move the binary somewhere on your `$PATH`:

```bash
mv pomo /usr/local/bin/
```

Or install directly with `go install`:

```bash
go install github.com/varuzhantrue/pomo@latest
```

## Usage

```
pomo start    # start a 25-minute focus session
pomo break    #  5-minute short break
pomo long     # 15-minute long break
pomo status   # show today's sessions and total focus time
```

**During a session**, press `Ctrl+C` to interrupt. You'll be asked whether to save the partial session.

**Example output:**

```
🍅 Focus Session #3  [████████████░░░░░░░░]  11:42 remaining

✅ Done! Take a break.
Today: 🍅🍅🍅  (3 sessions = 1h 15m focused)
```

## Data

Sessions are stored in `~/.pomo/sessions.json`. Each entry records the session type, start time, duration, and whether it was partial. The `status` command and session numbering are derived from this file.

## License

MIT

---

*Built with the assistance of Claude Code*

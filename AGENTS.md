# AGENTS

## Purpose

`radiofx` is a lightweight Go CLI for turning clean voice recordings into radio-style speech for simracing spotter / engineer voice packs.
It wraps `ffmpeg`, keeps the workflow offline and simple, and supports both CLI and a minimal terminal UI.

## Read First

When starting a new session, read in this order:

1. `AGENTS.md`
2. `.ai-memory/todo-plan.md` if it exists
3. `cmd/radiofx/main.go`
4. `internal/audio/processor.go`
5. `internal/config/config.go`
6. `presets.json`
7. `README.md`

## Project Facts

- `filter` mode processes voice only
- `noise` mode builds `start noise -> processed voice -> end noise`
- presets live in `presets.json`
- output is written next to the source file by default
- Windows users rely on the wrapper scripts in `windows/`

## Audio Chain

Voice processing is assembled in `internal/audio/processor.go`.
Current preset-driven stages:

- silence trim at start/end
- high-pass / low-pass EQ
- optional mid boost
- optional `compand`
- optional bit crusher
- volume adjustment
- loudness normalization
- limiter
- fade in / fade out
- start/end padding
- final sample-rate / mono formatting

## Key Files

- `cmd/radiofx/main.go`: CLI commands and argument flow
- `internal/audio/processor.go`: `ffmpeg` command assembly and filter chains
- `internal/config/config.go`: preset schema, defaults, validation
- `internal/ui/ui.go`: dependency-free interactive flow
- `presets.json`: shipping presets
- `windows/`: drag-and-drop wrappers
- `README.md`: user-facing behavior and usage

## Do Not Break

- Keep the project lightweight. Do not turn it into a DAW or GUI-heavy tool.
- Preserve the simple offline `ffmpeg` wrapper model.
- Preserve existing `filter` and `noise` workflows unless the change explicitly redesigns them.
- Do not break Windows drag-and-drop wrappers.
- Prefer preset-driven sound design over scattered one-off CLI flags.
- Keep backward compatibility for old `presets.json` files when reasonably possible.
- If you change DSP behavior, also update config validation and tests.
- If you change user-visible behavior or preset keys, update `README.md`.

## Change Strategy

- Read current behavior before editing; do not assume audio-chain order.
- Prefer extending `Preset` and the filter builders instead of adding ad hoc conditionals in multiple places.
- Keep new behavior deterministic unless randomness is an explicit feature.
- Add small focused tests around filter-string generation and config loading/validation.

## Verification

For Go code changes:

- run `gofmt -w` on changed Go files
- run `go test ./...`
- run `go build ./cmd/radiofx`

For behavior/config changes:

- confirm `presets.json` still loads
- confirm README matches the shipped preset keys and commands

## Local AI Memory

Local agent notes live in `.ai-memory/`.
That folder is intentionally gitignored and may contain session notes, plans, and future TODO items.
Read it before larger changes, but do not commit it.

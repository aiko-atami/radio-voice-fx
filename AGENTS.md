# AGENTS

## Project Summary

`radiofx` is a small Go CLI for turning clean voice recordings into radio-style speech for simracing spotter / engineer voice packs.
It wraps `ffmpeg`, keeps the workflow simple, and supports both CLI and a minimal terminal UI.

## Core Workflow

- `filter` mode processes the voice only
- `noise` mode builds `start noise -> processed voice -> end noise`
- Presets live in `presets.json`
- Output is written next to the source file with preset-based suffixes

## Current Audio Chain

Voice processing is assembled in `internal/audio/processor.go`.
Current preset-controlled stages:

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

- `cmd/radiofx/main.go`: CLI entrypoints and commands
- `internal/audio/processor.go`: `ffmpeg` command and filter-chain assembly
- `internal/config/config.go`: preset schema, defaults, validation
- `internal/ui/ui.go`: minimal interactive terminal flow
- `presets.json`: default shipping presets
- `windows/`: drag-and-drop wrappers for Windows users
- `README.md`: user-facing usage and packaging notes

## Testing

- Run `go test ./...`
- Run `go build ./cmd/radiofx`

There are unit tests for filter-chain generation and config defaults/validation.

## Local AI Memory

Local agent notes live in `.ai-memory/`.
That folder is intentionally gitignored and may contain session notes, plans, and future TODO items. If it exists, read it before making larger changes.

# Radio Voice FX

Cross-platform `ffmpeg` wrapper for turning clean voice recordings into radio-style speech, with optional random intro/outro noise.

## Download

Release builds are published automatically when you push a tag like `v0.1.0`.

Stable download links:

- [Windows x64 ZIP](https://github.com/aiko-atami/radio-voice-fx/releases/latest/download/radiofx-windows-x64.zip)
- [Linux x64 tar.gz](https://github.com/aiko-atami/radio-voice-fx/releases/latest/download/radiofx-linux-x64.tar.gz)

Create a release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub will mark the newest non-prerelease tag as the latest release, so `releases/latest/download/...` stays stable and is better than moving a tag like `last`.

## Requirements

- `ffmpeg` must be installed and available in `PATH`
- Go is only required if you want to build `radiofx` from source

Check `ffmpeg`:

```bash
ffmpeg -version
```

## Usage

Start the minimal terminal UI:

```bash
./radiofx
```

List presets:

```bash
./radiofx list-presets
```

Apply one preset in non-interactive mode:

```bash
./radiofx apply -preset clean_radio -mode filter voice.wav
./radiofx apply -preset harsh_digital -mode noise voice.wav
```

The tool writes output next to the source file.

## Windows Drag-And-Drop

If you want the old Explorer workflow, use the batch wrappers from `windows/`.
They look for `radiofx.exe` in the same directory first, then in the parent directory.

Available wrappers:

- `windows/clean-radio-filter.bat`
- `windows/walkie-talkie-filter.bat`
- `windows/harsh-digital-filter.bat`
- `windows/noisy-transmission-filter.bat`
- `windows/broken-signal-filter.bat`
- `windows/clean-radio-noise.bat`
- `windows/walkie-talkie-noise.bat`
- `windows/harsh-digital-noise.bat`
- `windows/noisy-transmission-noise.bat`
- `windows/broken-signal-noise.bat`

These wrappers are thin launchers around `radiofx.exe apply ... %*`.
They do not duplicate the filter logic.

Example output names when using drag-and-drop:

- `voice.wav` on `windows/clean-radio-filter.bat` -> `voice_clean_radio.wav`
- `voice.wav` on `windows/walkie-talkie-filter.bat` -> `voice_walkie.wav`
- `voice.wav` on `windows/harsh-digital-noise.bat` -> `voice_harsh_noise.wav`

## Presets

Presets live in `presets.json`.

Current defaults:

- `clean_radio` - light, clear radio voice
- `walkie_talkie` - classic handheld comms tone
- `harsh_digital` - stronger digital distortion with bit crushing
- `noisy_transmission` - more static and lower fidelity
- `broken_signal` - very narrow, damaged-sounding communication

Each preset defines voice-processing parameters such as:

- high-pass and low-pass cutoff
- output sample rate
- mono/stereo format
- compression
- bit crushing
- mid boost
- noise volume
- output filename suffixes

## Interactive Flow

The built-in terminal UI is intentionally minimal and dependency-free:

1. Choose mode: `filter` or `noise`
2. Choose preset from `presets.json`
3. Provide one or more input files
4. Let `radiofx` invoke `ffmpeg`

If files are passed as CLI arguments, the UI uses them directly:

```bash
./radiofx tui voice.wav other.wav
```

## Modes

### `filter`

Processes the voice only.

Example output:

```text
voice.wav -> voice_filtered.wav
```

### `noise`

Picks one random WAV from `start/` and one random WAV from `end/`, then concatenates:

1. start noise
2. processed voice
3. end noise

Example output:

```text
voice.wav -> voice_harsh_digital_w_noise.wav
```

## Notes

- `start/` and `end/` must contain `.wav` files for `noise` mode
- output files are never written over the original input
- `presets.json` is searched in the current directory first, then next to the executable
- wrapper `.bat` files in `windows/` work with `radiofx.exe` either next to the wrapper or in the parent directory
- old hardcoded scripts were moved to `legacy/`

## Build

The repository contains:

- `cmd/radiofx/` - Go CLI and minimal terminal UI
- `presets.json` - editable preset definitions
- `start/` - WAV files used as opening noise
- `end/` - WAV files used as ending noise
- `windows/` - Windows drag-and-drop wrappers for presets
- `legacy/` - old hardcoded batch scripts

Build for the current platform:

```bash
go build ./cmd/radiofx
```

Cross-compile examples:

```bash
GOOS=windows GOARCH=amd64 go build -o dist/radiofx.exe ./cmd/radiofx
GOOS=linux GOARCH=amd64 go build -o dist/radiofx-linux ./cmd/radiofx
```

GitHub Actions workflow:

- `.github/workflows/build-cli.yml` builds Windows and Linux artifacts
- pushing a tag matching `v*` also publishes release assets to GitHub Releases
- Windows artifacts include `radiofx.exe`, `presets.json`, `start/`, `end/`, and the wrapper `.bat` files in one directory

## Troubleshooting

### `ffmpeg` is not recognized

- install `ffmpeg`
- add it to `PATH`
- reopen the terminal

### `noise` mode says no WAV files were found

- add `.wav` clips into both `start/` and `end/`

### Processing fails on one file

- run the same command in a terminal to inspect the `ffmpeg` error output
- verify the input file is a readable audio format

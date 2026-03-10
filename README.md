# Radio Voice FX

Small Windows drag-and-drop tools for turning clean voice recordings into a radio-style sound and adding noise at the beginning and end.

This folder contains:

- `radio-filter.bat` - applies a radio-style voice filter.
- `add_noise.bat` - adds random intro/outro noise and also applies a stronger radio effect.
- `start/` - WAV files used as the opening noise.
- `end/` - WAV files used as the ending noise.

## What This Tool Does

These scripts are meant for quick offline processing of voice recordings with `ffmpeg`.

Typical workflow:

1. Take your clean audio recording.
2. Drag and drop it onto `radio-filter.bat` to get a distorted radio-style voice.
3. Drag and drop it onto `add_noise.bat` to get extra noise at the beginning and the end.

You can drop one file or multiple files at once onto either script.

## Requirement

You must have `ffmpeg` installed and available in your system `PATH`.

To check that it is installed, open Command Prompt and run:

```bat
ffmpeg -version
```

If Windows says the command is not recognized, install `ffmpeg` first and make sure the `ffmpeg.exe` location is added to `PATH`.

## How To Use

### Option 1: Create a Radio-Style Voice

Use `radio-filter.bat`.

Steps:

1. Prepare one or more audio files.
2. Select the files in Explorer.
3. Drag them onto `radio-filter.bat`.
4. Wait until processing finishes.

Output:

- A new WAV file is created next to each source file.
- The output name is:

```text
original_name_filtered.wav
```

Example:

```text
voice.wav -> voice_filtered.wav
```

What the script does:

- cuts low frequencies with a high-pass filter at 300 Hz
- cuts high frequencies with a low-pass filter at 3000 Hz
- reduces the sample rate to 11025 Hz
- applies dynamic compression
- slightly boosts volume

Result:

- narrower frequency range
- more compressed and less natural tone
- closer to a walkie-talkie / radio / transmission sound

### Option 2: Add Noise at the Beginning and End

Use `add_noise.bat`.

Steps:

1. Prepare one or more audio files.
2. Drag them onto `add_noise.bat`.
3. The script randomly picks one WAV from `start/` and one WAV from `end/`.
4. It places the start noise before the voice and the end noise after the voice.

Output:

- A new WAV file is created next to each source file.
- The output name is:

```text
original_name_w_noise.wav
```

Example:

```text
voice.wav -> voice_w_noise.wav
```

What this script does:

- adds a random opening noise from the `start/` folder
- adds a random closing noise from the `end/` folder
- lowers the noise volume to `0.3`
- converts audio to mono / 44100 Hz for the final mix
- applies radio-style EQ to the voice
- applies bit crushing for a rougher, dirtier transmission sound

Important note:

`add_noise.bat` does not only add noise. It also processes the main voice so the result sounds more like a noisy radio transmission.

## Batch Processing

Both scripts support multiple files.

That means you can select several recordings and drag all of them onto the `.bat` file in one action. The script will process them one by one.

## File Structure

```text
radio-voice-fx/
├─ radio-filter.bat
├─ add_noise.bat
├─ start/
│  └─ *.wav
└─ end/
   └─ *.wav
```

## Notes

- The noise clips inside `start/` and `end/` should stay in WAV format.
- If one of these folders is empty, `add_noise.bat` will stop with an error.
- The processed files are saved in the same folder as the original source file.
- Original files are not overwritten by default because each script creates a new output filename.

## Recommended Workflow

If you want a simple radio-style voice:

- use `radio-filter.bat`

If you want radio voice with static/noise before and after:

- use `add_noise.bat`

For fast usage, the simplest explanation is:

- drag your recording onto `radio-filter.bat` to get a radio-like distorted voice
- drag your recording onto `add_noise.bat` to get noise at the beginning and end

## Troubleshooting

### `ffmpeg` is not recognized

Cause:

- `ffmpeg` is not installed or not added to `PATH`

Fix:

- install `ffmpeg`
- add it to Windows `PATH`
- reopen Command Prompt or Explorer and try again

### No output file was created

Cause:

- unsupported input file
- `ffmpeg` failed during conversion

Fix:

- try processing the file from Command Prompt to see the error message
- make sure the input file is a valid audio file

### `add_noise.bat` says no WAV files were found

Cause:

- `start/` or `end/` is empty

Fix:

- place WAV noise files inside both folders

## Summary

This is a simple drag-and-drop Windows toolkit for voice post-processing:

- `radio-filter.bat` creates a filtered radio voice
- `add_noise.bat` creates a more aggressive radio effect with random noise at the start and end
- `ffmpeg` is required
- multiple files are supported

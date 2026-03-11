package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"radiofx/internal/audio"
	"radiofx/internal/config"
)

type Session struct {
	Mode   audio.Mode
	Preset config.Preset
	Files  []string
}

func RunInteractive(cfg *config.Config, initialFiles []string, in io.Reader, out io.Writer) (Session, error) {
	reader := bufio.NewReader(in)
	printHeader(out, cfg)

	mode, err := chooseMode(reader, out)
	if err != nil {
		return Session{}, err
	}

	preset, err := choosePreset(reader, out, cfg.Presets)
	if err != nil {
		return Session{}, err
	}

	files := initialFiles
	if len(files) == 0 {
		files, err = promptFiles(reader, out)
		if err != nil {
			return Session{}, err
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Mode:   %s\n", mode)
	fmt.Fprintf(out, "Preset: %s\n", preset.Name)
	fmt.Fprintf(out, "Files:  %d\n", len(files))

	if err := promptEnter(reader, out, "Press Enter to start processing..."); err != nil {
		return Session{}, err
	}

	return Session{
		Mode:   mode,
		Preset: preset,
		Files:  files,
	}, nil
}

func printHeader(out io.Writer, cfg *config.Config) {
	clearScreen(out)
	fmt.Fprintln(out, "radiofx")
	fmt.Fprintln(out, "Minimal terminal UI for radio-style voice processing")
	fmt.Fprintf(out, "Loaded presets: %d\n\n", len(cfg.Presets))
}

func chooseMode(reader *bufio.Reader, out io.Writer) (audio.Mode, error) {
	fmt.Fprintln(out, "Choose mode:")
	fmt.Fprintln(out, "  1) filter")
	fmt.Fprintln(out, "  2) noise")
	fmt.Fprintln(out)

	for {
		value, err := prompt(reader, out, "Mode [1-2]: ")
		if err != nil {
			return "", err
		}

		switch strings.TrimSpace(value) {
		case "1", "filter":
			return audio.ModeFilter, nil
		case "2", "noise":
			return audio.ModeNoise, nil
		default:
			fmt.Fprintln(out, "Enter 1 or 2.")
		}
	}
}

func choosePreset(reader *bufio.Reader, out io.Writer, presets []config.Preset) (config.Preset, error) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Choose preset:")
	for idx, preset := range presets {
		fmt.Fprintf(out, "  %d) %s - %s\n", idx+1, preset.Name, preset.Description)
	}
	fmt.Fprintln(out)

	for {
		value, err := prompt(reader, out, fmt.Sprintf("Preset [1-%d]: ", len(presets)))
		if err != nil {
			return config.Preset{}, err
		}

		index := parseIndex(value, len(presets))
		if index >= 0 {
			return presets[index], nil
		}

		for _, preset := range presets {
			if preset.Name == strings.TrimSpace(value) {
				return preset, nil
			}
		}

		fmt.Fprintln(out, "Enter a preset number or preset name.")
	}
}

func promptFiles(reader *bufio.Reader, out io.Writer) ([]string, error) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Enter one input file per line. Submit an empty line to finish.")

	files := []string{}
	for {
		value, err := prompt(reader, out, "File: ")
		if err != nil {
			return nil, err
		}
		value = strings.TrimSpace(value)
		if value == "" {
			if len(files) == 0 {
				fmt.Fprintln(out, "Add at least one file.")
				continue
			}
			return files, nil
		}
		files = append(files, value)
	}
}

func prompt(reader *bufio.Reader, out io.Writer, label string) (string, error) {
	fmt.Fprint(out, label)
	value, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return strings.TrimSpace(value), nil
		}
		return "", err
	}
	return strings.TrimSpace(value), nil
}

func promptEnter(reader *bufio.Reader, out io.Writer, label string) error {
	fmt.Fprint(out, label)
	_, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil
	}
	return err
}

func parseIndex(value string, count int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return -1
	}

	var index int
	if _, err := fmt.Sscanf(value, "%d", &index); err != nil {
		return -1
	}

	index--
	if index < 0 || index >= count {
		return -1
	}
	return index
}

func clearScreen(out io.Writer) {
	if file, ok := out.(*os.File); ok && file == os.Stdout {
		fmt.Fprint(out, "\033[H\033[2J")
	}
}

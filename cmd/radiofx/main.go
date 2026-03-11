package main

import (
	"flag"
	"fmt"
	"os"

	"radiofx/internal/audio"
	"radiofx/internal/config"
	"radiofx/internal/ui"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return runTUI(nil, "")
	}

	switch args[0] {
	case "tui":
		return runTUI(args[1:], "")
	case "apply":
		return runApply(args[1:])
	case "list-presets":
		return runListPresets(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return runTUI(args, "")
	}
}

func runTUI(args []string, configPath string) error {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	configFile := fs.String("config", configPath, "Path to presets.json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		return err
	}

	session, err := ui.RunInteractive(cfg, fs.Args(), os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	processor := audio.NewProcessor(cfg)
	return processor.ProcessFiles(session.Mode, session.Preset, "", session.Files)
}

func runApply(args []string) error {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	configPath := fs.String("config", "", "Path to presets.json")
	presetName := fs.String("preset", "", "Preset name")
	modeName := fs.String("mode", "", "Processing mode: filter or noise")
	suffix := fs.String("suffix", "", "Optional output suffix override")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *presetName == "" {
		return fmt.Errorf("apply: -preset is required")
	}
	if *modeName == "" {
		return fmt.Errorf("apply: -mode is required")
	}
	if len(fs.Args()) == 0 {
		return fmt.Errorf("apply: provide at least one input file")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}
	preset, err := cfg.PresetByName(*presetName)
	if err != nil {
		return err
	}
	mode, err := audio.ParseMode(*modeName)
	if err != nil {
		return err
	}

	processor := audio.NewProcessor(cfg)
	return processor.ProcessFiles(mode, preset, *suffix, fs.Args())
}

func runListPresets(args []string) error {
	fs := flag.NewFlagSet("list-presets", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	configPath := fs.String("config", "", "Path to presets.json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	for _, preset := range cfg.Presets {
		fmt.Printf("%s\t%s\n", preset.Name, preset.Description)
	}
	return nil
}

func printUsage() {
	fmt.Println("radiofx")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  radiofx                  Start interactive terminal UI")
	fmt.Println("  radiofx tui [files...]   Start interactive terminal UI")
	fmt.Println("  radiofx list-presets     Print configured presets")
	fmt.Println("  radiofx apply -preset <name> -mode <filter|noise> [files...]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  radiofx")
	fmt.Println("  radiofx tui voice.wav")
	fmt.Println("  radiofx apply -preset clean_radio -mode filter voice.wav")
	fmt.Println("  radiofx apply -preset harsh_digital -mode noise voice.wav")
}

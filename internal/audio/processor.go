package audio

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"radiofx/internal/config"
)

type Mode string

const (
	ModeFilter Mode = "filter"
	ModeNoise  Mode = "noise"
)

type Processor struct {
	cfg *config.Config
	rng *rand.Rand
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func ParseMode(raw string) (Mode, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(ModeFilter):
		return ModeFilter, nil
	case string(ModeNoise):
		return ModeNoise, nil
	default:
		return "", fmt.Errorf("unsupported mode %q", raw)
	}
}

func (p *Processor) ProcessFiles(mode Mode, preset config.Preset, suffix string, files []string) error {
	if len(files) == 0 {
		return errors.New("no input files provided")
	}

	for _, input := range files {
		if err := p.ProcessFile(mode, preset, suffix, input); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) ProcessFile(mode Mode, preset config.Preset, suffix, input string) error {
	inputPath, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("resolve input path: %w", err)
	}
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file %q: %w", inputPath, err)
	}

	outputPath := buildOutputPath(inputPath, extensionForPreset(preset), suffixForMode(preset, mode, suffix))
	cmd, err := p.buildCommand(mode, preset, inputPath, outputPath)
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("\n[%s] %s -> %s\n", strings.ToUpper(string(mode)), filepath.Base(inputPath), filepath.Base(outputPath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed for %q: %w", inputPath, err)
	}

	return nil
}

func (p *Processor) buildCommand(mode Mode, preset config.Preset, inputPath, outputPath string) (*exec.Cmd, error) {
	switch mode {
	case ModeFilter:
		return exec.Command("ffmpeg", "-y", "-i", inputPath, "-af", buildVoiceFilter(preset), outputPath), nil
	case ModeNoise:
		startNoise, err := randomAudioFile(p.cfg.StartPath(), p.rng)
		if err != nil {
			return nil, fmt.Errorf("pick start noise: %w", err)
		}
		endNoise, err := randomAudioFile(p.cfg.EndPath(), p.rng)
		if err != nil {
			return nil, fmt.Errorf("pick end noise: %w", err)
		}

		filterComplex := fmt.Sprintf(
			"[0:a]%s[a0];[1:a]%s[a1];[2:a]%s[a2];[a0][a1][a2]concat=n=3:v=0:a=1[out]",
			buildNoiseFilter(preset),
			buildVoiceFilterForNoise(preset),
			buildNoiseFilter(preset),
		)

		return exec.Command(
			"ffmpeg",
			"-y",
			"-i", startNoise,
			"-i", inputPath,
			"-i", endNoise,
			"-filter_complex", filterComplex,
			"-map", "[out]",
			outputPath,
		), nil
	default:
		return nil, fmt.Errorf("unsupported mode %q", mode)
	}
}

func buildVoiceFilter(preset config.Preset) string {
	return buildVoiceFilterWithFormat(preset, preset.SampleRate, preset.Mono)
}

func buildVoiceFilterForNoise(preset config.Preset) string {
	return buildVoiceFilterWithFormat(preset, preset.NoiseSampleRate, preset.NoiseMono)
}

func buildVoiceFilterWithFormat(preset config.Preset, sampleRate int, mono bool) string {
	parts := []string{}
	parts = append(parts, buildTrimFilters(preset)...)
	parts = append(parts,
		"highpass=f="+strconv.Itoa(preset.HighpassHz),
		"lowpass=f="+strconv.Itoa(preset.LowpassHz),
	)

	if preset.MidBoostHz > 0 && preset.MidBoostWidthHz > 0 && preset.MidBoostGainDB != 0 {
		parts = append(parts, fmt.Sprintf(
			"equalizer=f=%d:t=h:w=%d:g=%s",
			preset.MidBoostHz,
			preset.MidBoostWidthHz,
			floatString(preset.MidBoostGainDB),
		))
	}
	if preset.Compand {
		parts = append(parts, "compand")
	}
	if preset.CrusherBits > 0 {
		crusherParts := []string{"bits=" + strconv.Itoa(preset.CrusherBits)}
		if preset.CrusherMode != "" {
			crusherParts = append(crusherParts, "mode="+preset.CrusherMode)
		}
		if preset.CrusherAA > 0 {
			crusherParts = append(crusherParts, "aa="+strconv.Itoa(preset.CrusherAA))
		}
		parts = append(parts, "acrusher="+strings.Join(crusherParts, ":"))
	}

	if preset.Volume != 1 {
		parts = append(parts, "volume="+floatString(preset.Volume))
	}

	parts = append(parts, buildNormalizationFilters(preset)...)
	parts = append(parts, buildFadeFilters(preset)...)
	parts = append(parts, buildPaddingFilters(preset)...)
	parts = append(parts, buildFormatFilter(sampleRate, mono))

	return strings.Join(parts, ",")
}

func buildTrimFilters(preset config.Preset) []string {
	if !preset.TrimSilence {
		return nil
	}

	trim := fmt.Sprintf(
		"silenceremove=start_periods=1:start_duration=%s:start_threshold=%s:start_silence=0",
		durationStringFromMilliseconds(preset.TrimDurationMS),
		decibelString(preset.TrimThresholdDB),
	)

	// Reverse twice to apply the same leading-silence trim to the clip tail.
	return []string{
		trim,
		"areverse",
		trim,
		"areverse",
	}
}

func buildNormalizationFilters(preset config.Preset) []string {
	parts := []string{}

	if preset.NormalizeLUFS != 0 {
		parts = append(parts, fmt.Sprintf(
			"loudnorm=I=%s:LRA=%s:TP=%s:linear=true",
			floatString(preset.NormalizeLUFS),
			floatString(preset.NormalizeLRA),
			floatString(preset.NormalizeTP),
		))
	}
	if preset.LimiterCeilingDB != 0 {
		parts = append(parts, fmt.Sprintf(
			"alimiter=limit=%s:attack=%s:release=%s:level=true:latency=true",
			floatString(linearAmplitudeFromDB(preset.LimiterCeilingDB)),
			floatString(preset.LimiterAttackMS),
			floatString(preset.LimiterReleaseMS),
		))
	}

	return parts
}

func buildFadeFilters(preset config.Preset) []string {
	parts := []string{}
	if preset.FadeInMS > 0 {
		parts = append(parts, fmt.Sprintf(
			"afade=t=in:st=0:d=%s",
			durationStringFromMilliseconds(preset.FadeInMS),
		))
	}
	if preset.FadeOutMS > 0 {
		parts = append(parts,
			"areverse",
			fmt.Sprintf("afade=t=in:st=0:d=%s", durationStringFromMilliseconds(preset.FadeOutMS)),
			"areverse",
		)
	}
	return parts
}

func buildPaddingFilters(preset config.Preset) []string {
	parts := []string{}
	if preset.PadStartMS > 0 {
		parts = append(parts, fmt.Sprintf("adelay=delays=%d:all=1", preset.PadStartMS))
	}
	if preset.PadEndMS > 0 {
		parts = append(parts, fmt.Sprintf("apad=pad_dur=%s", durationStringFromMilliseconds(preset.PadEndMS)))
	}
	return parts
}

func buildNoiseFilter(preset config.Preset) string {
	parts := []string{}
	if preset.NoiseVolume > 0 {
		parts = append(parts, "volume="+floatString(preset.NoiseVolume))
	}
	parts = append(parts, buildFormatFilter(preset.NoiseSampleRate, preset.NoiseMono))
	return strings.Join(parts, ",")
}

func buildFormatFilter(sampleRate int, mono bool) string {
	if mono {
		return fmt.Sprintf("aformat=sample_rates=%d:channel_layouts=mono", sampleRate)
	}
	return fmt.Sprintf("aformat=sample_rates=%d", sampleRate)
}

func randomAudioFile(dir string, rng *rand.Rand) (string, error) {
	pattern := filepath.Join(dir, "*.wav")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no wav files found in %s", dir)
	}
	return matches[rng.Intn(len(matches))], nil
}

func buildOutputPath(inputPath, formatName, suffix string) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return base + "_" + suffix + "." + formatName
}

func extensionForPreset(preset config.Preset) string {
	if preset.OutputFormat == "" {
		return "wav"
	}
	return strings.TrimPrefix(preset.OutputFormat, ".")
}

func suffixForMode(preset config.Preset, mode Mode, override string) string {
	if override != "" {
		return override
	}
	switch mode {
	case ModeFilter:
		return preset.FilterSuffix
	case ModeNoise:
		return preset.NoiseSuffix
	default:
		return preset.Name
	}
}

func floatString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func durationStringFromMilliseconds(milliseconds int) string {
	return floatString(float64(milliseconds) / 1000)
}

func decibelString(value float64) string {
	return floatString(value) + "dB"
}

func linearAmplitudeFromDB(value float64) float64 {
	return math.Pow(10, value/20)
}

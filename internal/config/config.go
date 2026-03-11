package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultConfigName = "presets.json"

type Config struct {
	StartDir string   `json:"start_dir"`
	EndDir   string   `json:"end_dir"`
	Presets  []Preset `json:"presets"`

	baseDir string
}

type Preset struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	HighpassHz       int     `json:"highpass_hz"`
	LowpassHz        int     `json:"lowpass_hz"`
	SampleRate       int     `json:"sample_rate"`
	Mono             bool    `json:"mono"`
	Compand          bool    `json:"compand"`
	Volume           float64 `json:"volume"`
	CrusherBits      int     `json:"crusher_bits"`
	CrusherMode      string  `json:"crusher_mode"`
	CrusherAA        int     `json:"crusher_aa"`
	MidBoostHz       int     `json:"mid_boost_hz"`
	MidBoostWidthHz  int     `json:"mid_boost_width_hz"`
	MidBoostGainDB   float64 `json:"mid_boost_gain_db"`
	NormalizeLUFS    float64 `json:"normalize_lufs"`
	NormalizeLRA     float64 `json:"normalize_lra"`
	NormalizeTP      float64 `json:"normalize_tp"`
	LimiterCeilingDB float64 `json:"limiter_ceiling_db"`
	LimiterAttackMS  float64 `json:"limiter_attack_ms"`
	LimiterReleaseMS float64 `json:"limiter_release_ms"`
	TrimSilence      bool    `json:"trim_silence"`
	TrimThresholdDB  float64 `json:"trim_threshold_db"`
	TrimDurationMS   int     `json:"trim_duration_ms"`
	PadStartMS       int     `json:"pad_start_ms"`
	PadEndMS         int     `json:"pad_end_ms"`
	FadeInMS         int     `json:"fade_in_ms"`
	FadeOutMS        int     `json:"fade_out_ms"`
	NoiseVolume      float64 `json:"noise_volume"`
	NoiseSampleRate  int     `json:"noise_sample_rate"`
	NoiseMono        bool    `json:"noise_mono"`
	OutputFormat     string  `json:"output_format"`
	FilterSuffix     string  `json:"filter_suffix"`
	NoiseSuffix      string  `json:"noise_suffix"`
}

func Load(path string) (*Config, error) {
	configPath, err := resolveConfigPath(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.baseDir = filepath.Dir(configPath)
	cfg.applyDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) PresetByName(name string) (Preset, error) {
	for _, preset := range c.Presets {
		if preset.Name == name {
			return preset, nil
		}
	}
	return Preset{}, fmt.Errorf("preset %q not found", name)
}

func (c *Config) StartPath() string {
	return c.resolveAssetPath(c.StartDir)
}

func (c *Config) EndPath() string {
	return c.resolveAssetPath(c.EndDir)
}

func (c *Config) resolveAssetPath(dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	return filepath.Join(c.baseDir, dir)
}

func (c *Config) applyDefaults() {
	if c.StartDir == "" {
		c.StartDir = "start"
	}
	if c.EndDir == "" {
		c.EndDir = "end"
	}

	for i := range c.Presets {
		if c.Presets[i].SampleRate == 0 {
			c.Presets[i].SampleRate = 11025
		}
		if c.Presets[i].NoiseSampleRate == 0 {
			c.Presets[i].NoiseSampleRate = 44100
		}
		if c.Presets[i].OutputFormat == "" {
			c.Presets[i].OutputFormat = "wav"
		}
		if c.Presets[i].FilterSuffix == "" {
			c.Presets[i].FilterSuffix = fmt.Sprintf("%s_filtered", c.Presets[i].Name)
		}
		if c.Presets[i].NoiseSuffix == "" {
			c.Presets[i].NoiseSuffix = fmt.Sprintf("%s_w_noise", c.Presets[i].Name)
		}
		if c.Presets[i].Volume == 0 {
			c.Presets[i].Volume = 1
		}
		if c.Presets[i].NormalizeLUFS != 0 {
			if c.Presets[i].NormalizeLRA == 0 {
				c.Presets[i].NormalizeLRA = 7
			}
			if c.Presets[i].NormalizeTP == 0 {
				c.Presets[i].NormalizeTP = -2
			}
		}
		if c.Presets[i].LimiterCeilingDB != 0 {
			if c.Presets[i].LimiterAttackMS == 0 {
				c.Presets[i].LimiterAttackMS = 5
			}
			if c.Presets[i].LimiterReleaseMS == 0 {
				c.Presets[i].LimiterReleaseMS = 50
			}
		}
		if c.Presets[i].TrimSilence {
			if c.Presets[i].TrimThresholdDB == 0 {
				c.Presets[i].TrimThresholdDB = -50
			}
			if c.Presets[i].TrimDurationMS == 0 {
				c.Presets[i].TrimDurationMS = 30
			}
		}
	}
}

func (c *Config) validate() error {
	if len(c.Presets) == 0 {
		return errors.New("config has no presets")
	}
	for _, preset := range c.Presets {
		if preset.Name == "" {
			return errors.New("preset name cannot be empty")
		}
		if preset.HighpassHz <= 0 || preset.LowpassHz <= 0 {
			return fmt.Errorf("preset %q must define highpass_hz and lowpass_hz", preset.Name)
		}
		if preset.SampleRate <= 0 {
			return fmt.Errorf("preset %q must define a positive sample_rate", preset.Name)
		}
		if preset.NormalizeLUFS != 0 {
			if preset.NormalizeLUFS >= 0 {
				return fmt.Errorf("preset %q must define normalize_lufs below 0", preset.Name)
			}
			if preset.NormalizeLRA <= 0 {
				return fmt.Errorf("preset %q must define a positive normalize_lra when normalize_lufs is set", preset.Name)
			}
			if preset.NormalizeTP >= 0 {
				return fmt.Errorf("preset %q must define normalize_tp below 0 when normalize_lufs is set", preset.Name)
			}
		}
		if preset.LimiterCeilingDB != 0 {
			if preset.LimiterCeilingDB >= 0 {
				return fmt.Errorf("preset %q must define limiter_ceiling_db below 0", preset.Name)
			}
			if preset.LimiterAttackMS <= 0 || preset.LimiterReleaseMS <= 0 {
				return fmt.Errorf("preset %q must define positive limiter attack/release values", preset.Name)
			}
		}
		if preset.TrimSilence {
			if preset.TrimThresholdDB >= 0 {
				return fmt.Errorf("preset %q must define trim_threshold_db below 0 when trim_silence is enabled", preset.Name)
			}
			if preset.TrimDurationMS <= 0 {
				return fmt.Errorf("preset %q must define a positive trim_duration_ms when trim_silence is enabled", preset.Name)
			}
		}
		if preset.PadStartMS < 0 || preset.PadEndMS < 0 {
			return fmt.Errorf("preset %q must define non-negative pad durations", preset.Name)
		}
		if preset.FadeInMS < 0 || preset.FadeOutMS < 0 {
			return fmt.Errorf("preset %q must define non-negative fade durations", preset.Name)
		}
	}
	return nil
}

func resolveConfigPath(path string) (string, error) {
	if path != "" {
		return filepath.Abs(path)
	}

	if cwdPath, err := filepath.Abs(DefaultConfigName); err == nil {
		if _, statErr := os.Stat(cwdPath); statErr == nil {
			return cwdPath, nil
		}
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable path: %w", err)
	}
	exeConfigPath := filepath.Join(filepath.Dir(exePath), DefaultConfigName)
	if _, err := os.Stat(exeConfigPath); err == nil {
		return exeConfigPath, nil
	}

	return "", fmt.Errorf("could not find %s in cwd or executable directory", DefaultConfigName)
}

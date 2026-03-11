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
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	HighpassHz      int     `json:"highpass_hz"`
	LowpassHz       int     `json:"lowpass_hz"`
	SampleRate      int     `json:"sample_rate"`
	Mono            bool    `json:"mono"`
	Compand         bool    `json:"compand"`
	Volume          float64 `json:"volume"`
	CrusherBits     int     `json:"crusher_bits"`
	CrusherMode     string  `json:"crusher_mode"`
	CrusherAA       int     `json:"crusher_aa"`
	MidBoostHz      int     `json:"mid_boost_hz"`
	MidBoostWidthHz int     `json:"mid_boost_width_hz"`
	MidBoostGainDB  float64 `json:"mid_boost_gain_db"`
	NoiseVolume     float64 `json:"noise_volume"`
	NoiseSampleRate int     `json:"noise_sample_rate"`
	NoiseMono       bool    `json:"noise_mono"`
	OutputFormat    string  `json:"output_format"`
	FilterSuffix    string  `json:"filter_suffix"`
	NoiseSuffix     string  `json:"noise_suffix"`
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

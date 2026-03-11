package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAppliesDefaultsForLoudnessLimiterAndTrim(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "presets.json")

	data := []byte(`{
  "presets": [
    {
      "name": "spotter",
      "highpass_hz": 280,
      "lowpass_hz": 3400,
      "normalize_lufs": -16,
      "limiter_ceiling_db": -2,
      "trim_silence": true
    }
  ]
}`)

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	preset := cfg.Presets[0]
	if preset.SampleRate != 11025 {
		t.Fatalf("expected default sample rate, got %d", preset.SampleRate)
	}
	if preset.NormalizeLRA != 7 {
		t.Fatalf("expected default normalize_lra, got %v", preset.NormalizeLRA)
	}
	if preset.NormalizeTP != -2 {
		t.Fatalf("expected default normalize_tp, got %v", preset.NormalizeTP)
	}
	if preset.LimiterAttackMS != 5 || preset.LimiterReleaseMS != 50 {
		t.Fatalf("expected default limiter timings, got attack=%v release=%v", preset.LimiterAttackMS, preset.LimiterReleaseMS)
	}
	if preset.TrimThresholdDB != -50 || preset.TrimDurationMS != 30 {
		t.Fatalf("expected default trim settings, got threshold=%v duration=%d", preset.TrimThresholdDB, preset.TrimDurationMS)
	}
}

func TestLoadRejectsInvalidTrimThreshold(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "presets.json")

	data := []byte(`{
  "presets": [
    {
      "name": "spotter",
      "highpass_hz": 280,
      "lowpass_hz": 3400,
      "trim_silence": true,
      "trim_threshold_db": 1
    }
  ]
}`)

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := Load(configPath); err == nil {
		t.Fatal("expected invalid trim threshold to fail validation")
	}
}

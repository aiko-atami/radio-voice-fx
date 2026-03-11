package audio

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"radiofx/internal/config"
)

func TestBuildVoiceFilterAddsLoudnessTrimAndTimingFilters(t *testing.T) {
	preset := config.Preset{
		HighpassHz:       300,
		LowpassHz:        3400,
		SampleRate:       12000,
		Mono:             true,
		Compand:          true,
		Volume:           1.2,
		CrusherBits:      8,
		CrusherMode:      "log",
		CrusherAA:        1,
		MidBoostHz:       1800,
		MidBoostWidthHz:  1200,
		MidBoostGainDB:   3,
		NormalizeLUFS:    -16,
		NormalizeLRA:     7,
		NormalizeTP:      -2,
		LimiterCeilingDB: -1,
		LimiterAttackMS:  4,
		LimiterReleaseMS: 60,
		TrimSilence:      true,
		TrimThresholdDB:  -48,
		TrimDurationMS:   25,
		PadStartMS:       40,
		PadEndMS:         80,
		FadeInMS:         15,
		FadeOutMS:        45,
	}

	filter := buildVoiceFilter(preset)

	trimFilter := "silenceremove=start_periods=1:start_duration=0.025:start_threshold=-48dB:start_silence=0"
	expected := []string{
		trimFilter,
		"areverse",
		trimFilter,
		"areverse",
		"highpass=f=300",
		"lowpass=f=3400",
		"equalizer=f=1800:t=h:w=1200:g=3",
		"compand",
		"acrusher=bits=8:mode=log:aa=1",
		"volume=1.2",
		"loudnorm=I=-16:LRA=7:TP=-2:linear=true",
		fmt.Sprintf(
			"alimiter=limit=%s:attack=4:release=60:level=true:latency=true",
			floatString(linearAmplitudeFromDB(-1)),
		),
		"afade=t=in:st=0:d=0.015",
		"areverse",
		"afade=t=in:st=0:d=0.045",
		"areverse",
		"adelay=delays=40:all=1",
		"apad=pad_dur=0.08",
		"aformat=sample_rates=12000:channel_layouts=mono",
	}

	if got := strings.Split(filter, ","); !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected filter chain:\nwant: %#v\ngot:  %#v", expected, got)
	}
}

func TestBuildVoiceFilterForNoiseUsesNoiseFormat(t *testing.T) {
	preset := config.Preset{
		HighpassHz:      320,
		LowpassHz:       3000,
		SampleRate:      11025,
		NoiseSampleRate: 44100,
		NoiseMono:       false,
	}

	filter := buildVoiceFilterForNoise(preset)
	if !strings.Contains(filter, "aformat=sample_rates=44100") {
		t.Fatalf("expected noise sample rate in filter, got %q", filter)
	}
	if strings.Contains(filter, "channel_layouts=mono") {
		t.Fatalf("expected noise voice filter to preserve non-mono layout, got %q", filter)
	}
}

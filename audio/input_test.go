package audio_test

import (
	"testing"

	"github.com/dieklingel/core/audio"
)

func TestNewAudioInputDeviceEmptyPipeline(t *testing.T) {
	input, err := audio.NewInput("")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if input != nil {
		t.Fatalf("expected nil, got %v", input)
	}
}

func TestNewAudioInputDeviceSuccesfull(t *testing.T) {
	input, err := audio.NewInput("audiotestsrc ! audio/x-raw, format=S16LE, layout=interleaved, rate=48000, channels=1 ! appsink name=rawsink")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if input == nil {
		t.Fatalf("expected not nil, got nil")
	}
}

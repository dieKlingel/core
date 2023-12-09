package audio_test

import (
	"testing"

	"github.com/dieklingel/core/audio"
)

func TestAudioCodecToPipelineElementsFromOpus(t *testing.T) {
	source, pipeline, sink := audio.OpusAudioCodec.ToPipelineElements()
	if source == nil {
		t.Fatal("source is nil")
	}
	if pipeline == nil {
		t.Fatal("pipeline is nil")
	}
	if sink == nil {
		t.Fatal("sink is nil")
	}
}

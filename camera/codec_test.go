package camera_test

import (
	"testing"

	"github.com/dieklingel/core/camera"
)

func TestAudioCodecToPipelineElementsFromX264(t *testing.T) {
	source, pipeline, sink := camera.X264CameraCodec.ToPipelineElements()
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

func TestAudioCodecToPipelineElementsFromMJPEG(t *testing.T) {
	source, pipeline, sink := camera.MJPEGCameraCodec.ToPipelineElements()
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

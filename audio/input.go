package audio

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type Input struct {
	sink     *app.Sink
	pipeline *gst.Pipeline
	streams  map[*AudioStream]struct{}
}

// NewInput creates a new audio input device from the given
// pipeline. The pipeline must have exactly one sink element. The sink element
// will be used to receive audio data from the device.
//
// example pipeline:
//
//	audiotestsrc ! audio/x-raw, format=S16LE, layout=interleaved, rate=48000, channels=1 ! appsink
func NewInput(definiton string) (*Input, error) {
	audioInputDevice := &Input{
		streams: make(map[*AudioStream]struct{}),
	}

	pipeline, err := gst.NewPipelineFromString(definiton)
	if err != nil {
		return nil, err
	}
	audioInputDevice.pipeline = pipeline

	sinks, err := pipeline.GetSinkElements()
	if err != nil {
		return nil, err
	}
	if len(sinks) != 1 {
		return nil, fmt.Errorf("audio input device must have exactly one sink")
	}
	audioInputDevice.sink = app.SinkFromElement(sinks[0])

	return audioInputDevice, nil
}

// Tee creates a new audio stream from the input device. The stream will be
// encoded with the given codec. The stream will be added to the input device
// and will receive audio data from the device.
func (input *Input) Tee(codec AudioCodec) *AudioStream {
	stream := NewAudioStreamFromInput(codec, input)
	input.streams[stream] = struct{}{}

	if len(input.streams) == 1 {
		input.pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
			switch msg.Type() {
			case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
				input.pipeline.BlockSetState(gst.StateNull)
			case gst.MessageError: // Error messages are always fatal
				err := msg.ParseError()
				fmt.Println("ERROR:", err.Error())
				if debug := err.DebugString(); debug != "" {
					fmt.Println("DEBUG:", debug)
				}
			}
			return true
		})
		input.sink.SetCallbacks(&app.SinkCallbacks{
			NewSampleFunc: func(s *app.Sink) gst.FlowReturn {
				sample := s.PullSample()
				if sample == nil {
					return gst.FlowEOS
				}

				for stream := range input.streams {
					stream.source.PushSample(sample)
				}

				return gst.FlowOK
			},
		})
		input.pipeline.SetState(gst.StatePlaying)
	}

	return stream
}

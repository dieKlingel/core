package audio

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type AudioStream struct {
	frame    chan gst.Sample
	finished chan bool
	codec    AudioCodec
	sink     *app.Sink
	source   *app.Source
	pipeline *gst.Pipeline
	input    *Input
}

func NewAudioStreamFromInput(codec AudioCodec, input *Input) *AudioStream {
	src, pipeline, sink := codec.ToPipelineElements()
	audioStream := &AudioStream{
		frame:    make(chan gst.Sample),
		finished: make(chan bool),
		codec:    codec,
		sink:     sink,
		source:   src,
		pipeline: pipeline,
		input:    input,
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			pipeline.BlockSetState(gst.StateNull)
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
		}
		return true
	})

	sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(s *app.Sink) gst.FlowReturn {
			sample := s.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}
			audioStream.frame <- *sample

			return gst.FlowOK
		},
	})

	pipeline.SetState(gst.StatePlaying)
	return audioStream
}

// Close closes the audio stream. The stream will be removed from the input
// device and will no longer receive audio data.
func (stream *AudioStream) Close() {
	delete(stream.input.streams, stream)
	close(stream.frame)
	close(stream.finished)
	err := stream.pipeline.BlockSetState(gst.StateNull)
	if err != nil {
		panic(err)
	}
	if len(stream.input.streams) == 0 {
		stream.input.pipeline.SetState(gst.StateNull)
	}
}

func (stream *AudioStream) FinishAndClose() {
	stream.finished <- true
	stream.Close()
}

// Codec returns the codec of the audio stream.
func (stream *AudioStream) Codec() AudioCodec {
	return stream.codec
}

// Frame returns a channel of audio frames.
func (stream *AudioStream) Frame() chan gst.Sample {
	return stream.frame
}

// Finished returns a channel that is closed when the stream is finished.
func (stream *AudioStream) Finished() chan bool {
	return stream.finished
}

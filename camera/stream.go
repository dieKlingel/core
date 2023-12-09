package camera

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type CameraStream struct {
	frame    chan *gst.Sample
	finished chan bool
	codec    CameraCodec
	sink     *app.Sink
	source   *app.Source
	pipeline *gst.Pipeline
	camera   *Camera
}

func NewStream(codec CameraCodec, camera *Camera) *CameraStream {
	src, pipeline, sink := codec.ToPipelineElements()
	cameraStream := &CameraStream{
		frame:    make(chan *gst.Sample),
		finished: make(chan bool),
		codec:    codec,
		sink:     sink,
		source:   src,
		pipeline: pipeline,
		camera:   camera,
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
			cameraStream.frame <- sample

			return gst.FlowOK
		},
	})

	pipeline.SetState(gst.StatePlaying)
	return cameraStream
}

// Close closes the camera stream. The stream will be removed from the input
// device and will no longer receive video data.
func (stream *CameraStream) Close() {
	delete(stream.camera.streams, stream)
	close(stream.frame)
	close(stream.finished)
	err := stream.pipeline.BlockSetState(gst.StateNull)
	if err != nil {
		panic(err)
	}
	if len(stream.camera.streams) == 0 {
		stream.camera.pipeline.SetState(gst.StateNull)
	}
}

func (stream *CameraStream) FinishAndClose() {
	stream.finished <- true
	stream.Close()
}

// Codec returns the codec of the camera stream.
func (stream *CameraStream) Codec() CameraCodec {
	return stream.codec
}

// Frame returns a channel of camera frames.
func (stream *CameraStream) Frame() chan *gst.Sample {
	return stream.frame
}

// Finished returns a channel that is closed when the stream is finished.
func (stream *CameraStream) Finished() chan bool {
	return stream.finished
}

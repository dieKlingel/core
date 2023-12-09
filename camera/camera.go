package camera

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

func init() {
	gst.Init(nil)
}

type Camera struct {
	sink     *app.Sink
	pipeline *gst.Pipeline
	streams  map[*CameraStream]struct{}
}

// New creates a new camera device from the given
// pipeline. The pipeline must have exactly one sink element. The sink element
// will be used to receive video data from the device.
//
// example pipeline:
//
//	videotestsrc ! video/x-raw, framerate=30/1, width=1280, height=720 ! appsink name=rawsink
func New(definiton string) (*Camera, error) {
	camera := &Camera{
		streams: make(map[*CameraStream]struct{}),
	}

	pipeline, err := gst.NewPipelineFromString(definiton)
	if err != nil {
		return nil, err
	}
	camera.pipeline = pipeline

	sinks, err := pipeline.GetSinkElements()
	if err != nil {
		return nil, err
	}
	if len(sinks) != 1 {
		return nil, fmt.Errorf("camera must have exactly one sink")
	}
	camera.sink = app.SinkFromElement(sinks[0])

	return camera, nil
}

func (camera *Camera) Tee(codec CameraCodec) *CameraStream {
	stream := NewStream(codec, camera)
	camera.streams[stream] = struct{}{}

	if len(camera.streams) == 1 {
		camera.pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
			switch msg.Type() {
			case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
				camera.pipeline.BlockSetState(gst.StateNull)
			case gst.MessageError: // Error messages are always fatal
				err := msg.ParseError()
				fmt.Println("ERROR:", err.Error())
				camera.pipeline.BlockSetState(gst.StateNull)
			}
			return true
		})
		camera.pipeline.SetState(gst.StatePlaying)

		camera.sink.SetCallbacks(&app.SinkCallbacks{
			NewSampleFunc: func(s *app.Sink) gst.FlowReturn {
				sample := s.PullSample()
				if sample == nil {
					return gst.FlowEOS
				}

				for stream := range camera.streams {
					stream.source.PushSample(sample)
				}

				return gst.FlowOK
			},
		})
		camera.pipeline.SetState(gst.StatePlaying)
	}

	return stream
}

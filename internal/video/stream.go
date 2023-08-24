package video

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

type Stream struct {
	Frame    chan gst.Sample
	Finished chan bool
	id       string
	pipeline *gst.Pipeline
	src      *app.Source
	sink     *app.Sink
	camera   *Camera
}

func NewStream(pipeline string) (*Stream, error) {
	pipe, err := gst.NewPipelineFromString(pipeline)
	if err != nil {
		return nil, err
	}

	pipe.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			pipe.BlockSetState(gst.StateNull)
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	stream := &Stream{
		Frame:    make(chan gst.Sample),
		Finished: make(chan bool),
		pipeline: pipe,
		id:       uuid.New().String(),
	}

	sink, err := pipe.GetElementByName("sink")
	if err != nil {
		return nil, err
	}
	stream.sink = app.SinkFromElement(sink)

	stream.sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := appSink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}
			stream.Frame <- *sample

			return gst.FlowOK
		},
	})

	src, err := pipe.GetElementByName("src")
	if err != nil {
		return nil, err
	}
	stream.src = app.SrcFromElement(src)

	stream.pipeline.SetState(gst.StatePlaying)
	return stream, nil
}

func (stream *Stream) ID() string {
	return stream.id
}

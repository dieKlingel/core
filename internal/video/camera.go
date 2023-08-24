package video

import (
	"fmt"
	"log"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

type Camera struct {
	streams  map[string]*Stream
	pipeline *gst.Pipeline
	loop     *glib.MainLoop
	sink     *app.Sink
}

func NewCamera(pipeline string) (*Camera, error) {
	pipe, err := gst.NewPipelineFromString(pipeline)
	if err != nil {
		return nil, err
	}

	if _, err := pipe.GetElementByName("rawsink"); err != nil {
		return nil, err
	}

	camera := &Camera{
		streams:  make(map[string]*Stream),
		pipeline: pipe,
	}

	return camera, nil
}

func (cam *Camera) AddStream(stream *Stream) {
	if _, exists := cam.streams[stream.ID()]; exists {
		caps := cam.sink.GetCaps()
		stream.src.SetCaps(caps)
		return
	}
	if stream.camera != nil {
		panic("a stream can only be added to one camera")
	}

	cam.streams[stream.ID()] = stream
	stream.camera = cam
	cam.openSafe()
}

func (cam *Camera) RemoveStream(stream *Stream) {
	delete(cam.streams, stream.ID())
	stream.camera = nil
	cam.closeSafe()
}

func (cam *Camera) openSafe() {
	if cam.loop != nil && cam.loop.IsRunning() {
		return
	}
	log.Println("open camera")

	cam.loop = glib.NewMainLoop(glib.MainContextDefault(), false)

	cam.pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			cam.pipeline.BlockSetState(gst.StateNull)
			cam.loop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			cam.loop.Quit()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	appsink, err := cam.pipeline.GetElementByName("rawsink")
	if err != nil {
		panic(err.Error())
	}
	cam.sink = app.SinkFromElement(appsink)

	cam.sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := cam.sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}

			for _, stream := range cam.streams {
				stream.src.PushSample(sample)
			}

			return gst.FlowOK
		},
	})

	cam.pipeline.SetState(gst.StatePlaying)
	go cam.loop.Run()
}

func (cam *Camera) closeSafe() {
	if len(cam.streams) != 0 || !cam.loop.IsRunning() {
		return
	}

	log.Println("close camera")

	cam.pipeline.BlockSetState(gst.StateNull)
	cam.loop.Quit()
}

package gmedia

import (
	"errors"
	"fmt"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
	"golang.org/x/exp/maps"
)

type VideoSrc struct {
	src      string
	tracks   map[string]*webrtc.TrackLocalStaticSample
	loop     *glib.MainLoop
	pipeline *gst.Pipeline
}

func NewVideoSrc(pipeline string) *VideoSrc {
	return &VideoSrc{
		src:    pipeline,
		tracks: make(map[string]*webrtc.TrackLocalStaticSample),
	}
}

func (src *VideoSrc) AddH264VideoTrack(track *webrtc.TrackLocalStaticSample) {
	src.tracks[track.ID()] = track
}

func (src *VideoSrc) RemoveH264VideoTrack(track *webrtc.TrackLocalStaticSample) {
	delete(src.tracks, track.ID())
}

func (src *VideoSrc) Open() error {
	if src.loop != nil {
		return errors.New("cannot open a VideoSrc, when it is already running")
	}

	loop := glib.NewMainLoop(glib.MainContextDefault(), false)
	pipeline, err := gst.NewPipelineFromString(src.src)
	if err != nil {
		return err
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			pipeline.BlockSetState(gst.StateNull)
			src.loop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			src.loop.Quit()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	h264sinkElement, err := pipeline.GetElementByName("h264sink")
	if err != nil {
		return err
	}
	h264sink := app.SinkFromElement(h264sinkElement)

	h264sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := h264sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}

			for _, track := range src.tracks {
				track.WriteSample(media.Sample{
					Data:     buffer.Bytes(),
					Duration: 1 * time.Millisecond, // use 1ms, because duration is incorrect when used with libcamerasrc, which is our preffered way
				})
			}

			return gst.FlowOK
		},
	})

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)
	src.pipeline = pipeline
	src.loop = loop

	go src.loop.Run()
	return nil
}

func (src *VideoSrc) Close() {
	if src.loop == nil {
		return
	}
	src.pipeline.BlockSetState(gst.StateNull)
	src.loop.Quit()
}

func (src *VideoSrc) Tracks() []*webrtc.TrackLocalStaticSample {
	return maps.Values(src.tracks)
}

func (src *VideoSrc) IsOpen() bool {
	return src.loop != nil
}

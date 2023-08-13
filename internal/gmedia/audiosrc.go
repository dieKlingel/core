package gmedia

import (
	"errors"
	"fmt"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
	"golang.org/x/exp/maps"
)

type AudioSrc struct {
	src      string
	tracks   map[string]*webrtc.TrackLocalStaticSample
	loop     *glib.MainLoop
	pipeline *gst.Pipeline
}

func NewAudioSrc(pipeline string) *AudioSrc {
	return &AudioSrc{
		src:    pipeline,
		tracks: make(map[string]*webrtc.TrackLocalStaticSample),
	}
}

func (src *AudioSrc) AddOpusAudioTrack(track *webrtc.TrackLocalStaticSample) {
	src.tracks[track.ID()] = track
}

func (src *AudioSrc) RemoveOpusAudioTrack(track *webrtc.TrackLocalStaticSample) {
	delete(src.tracks, track.ID())
}

func (src *AudioSrc) Open() error {
	if src.loop != nil {
		return errors.New("cannot open a Audio, when it is already running")
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

	opusSinkElement, err := pipeline.GetElementByName("opussink")
	if err != nil {
		return err
	}
	opusSink := app.SinkFromElement(opusSinkElement)

	opusSink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := opusSink.PullSample()
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
					Duration: buffer.Duration(),
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

func (src *AudioSrc) Close() {
	if src.loop == nil {
		return
	}
	src.pipeline.BlockSetState(gst.StateNull)
	src.loop.Quit()
}

func (src *AudioSrc) Tracks() []*webrtc.TrackLocalStaticSample {
	return maps.Values(src.tracks)
}

func (src *AudioSrc) IsOpen() bool {
	return src.loop != nil
}

package gmedia

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

type gaudio struct {
	src      string
	tracks   map[string]*webrtc.TrackLocalStaticSample
	mutex    sync.Mutex
	loop     *glib.MainLoop
	pipeline *gst.Pipeline
}

var audio = gaudio{
	tracks: make(map[string]*webrtc.TrackLocalStaticSample, 0),
}

func init() {
	gst.Init(nil)
}

func SetAudioSrc(src string) error {
	if len(audio.src) != 0 {
		return errors.New("cannot set audiosrc if already set")
	}
	audio.src = src
	if len(audio.tracks) > 0 {
		openAudio()
	}
	return nil
}

func AddAudioTrack(track *webrtc.TrackLocalStaticSample) {
	audio.mutex.Lock()
	audio.tracks[track.ID()] = track
	println(len(audio.src))
	if len(audio.tracks) == 1 && len(audio.src) != 0 {
		openAudio()
	}
	audio.mutex.Unlock()
}

func RemoveAudioTrack(track *webrtc.TrackLocalStaticSample) {
	audio.mutex.Lock()
	delete(audio.tracks, track.ID())
	if len(audio.tracks) == 0 {
		closeAudio()
	}
	audio.mutex.Unlock()
}

func openAudio() {
	audio.loop = glib.NewMainLoop(glib.MainContextDefault(), false)

	pipeline, err := gst.NewPipelineFromString(audio.src + " ! opusenc ! appsink name=sink1")
	if err != nil {
		panic(err.Error())
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			pipeline.BlockSetState(gst.StateNull)
			audio.loop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			audio.loop.Quit()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	si, er := pipeline.GetElementByName("sink1")
	if er != nil {
		panic(er.Error())
	}

	sink := app.SinkFromElement(si)
	if err != nil {
		panic(err.Error())
	}

	sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}

			for _, track := range audio.tracks {
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
	audio.pipeline = pipeline

	// Block and iterate on the main loop
	go audio.loop.Run()
}

func closeAudio() {
	if audio.pipeline == nil {
		return
	}
	audio.pipeline.SetState(gst.StateNull)
	audio.loop.Quit()
}

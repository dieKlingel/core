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

type gvideo struct {
	src      string
	tracks   map[string]*webrtc.TrackLocalStaticSample
	mutex    sync.Mutex
	loop     *glib.MainLoop
	pipeline *gst.Pipeline
}

var video = gvideo{
	tracks: make(map[string]*webrtc.TrackLocalStaticSample, 0),
}

func init() {
	gst.Init(nil)
}

func SetVideoSrc(src string) error {
	if len(video.src) != 0 {
		return errors.New("cannot set videosrc if already set")
	}
	video.src = src
	if len(video.tracks) > 0 {
		openVideo()
	}
	return nil
}

func AddVideoTrack(track *webrtc.TrackLocalStaticSample) {
	video.mutex.Lock()
	video.tracks[track.ID()] = track
	println(len(video.src))
	if len(video.tracks) == 1 && len(video.src) != 0 {
		openVideo()
	}
	video.mutex.Unlock()
}

func RemoveVideoTrack(track *webrtc.TrackLocalStaticSample) {
	video.mutex.Lock()
	delete(video.tracks, track.ID())

	if len(video.tracks) == 0 {
		closeVideo()
	}
	video.mutex.Unlock()
}

func openVideo() {
	video.loop = glib.NewMainLoop(glib.MainContextDefault(), false)

	pipeline, err := gst.NewPipelineFromString(video.src + " ! videoconvert ! vp8enc error-resilient=partitions keyframe-max-dist=10 auto-alt-ref=true cpu-used=5 deadline=1 ! appsink name=sink")
	if err != nil {
		panic(err.Error())
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			pipeline.BlockSetState(gst.StateNull)
			video.loop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			video.loop.Quit()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	si, er := pipeline.GetElementByName("sink")
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

			for _, track := range video.tracks {
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
	video.pipeline = pipeline

	// Block and iterate on the main loop
	go video.loop.Run()
}

func closeVideo() {

	video.pipeline.SetState(gst.StateNull)
	video.loop.Quit()
}

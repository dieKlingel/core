package gmedia

import (
	"errors"
	"fmt"

	"github.com/pion/webrtc/v3"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

type RemoteAudioSink struct {
	src      string
	track    *webrtc.TrackRemote
	loop     *glib.MainLoop
	pipeline *gst.Pipeline
}

func NewRemoteAudioSink(pipeline string, track *webrtc.TrackRemote) *RemoteAudioSink {
	return &RemoteAudioSink{
		src:   pipeline,
		track: track,
	}
}

func (sink *RemoteAudioSink) Open() error {
	if sink.loop != nil {
		return errors.New("cannot open a Audio, when it is already running")
	}

	loop := glib.NewMainLoop(glib.MainContextDefault(), false)
	pipeline, err := gst.NewPipelineFromString(sink.src)
	if err != nil {
		return err
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			print("EOS")
			pipeline.BlockSetState(gst.StateNull)
			sink.loop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			sink.loop.Quit()
		default:
			//println(msg.String())
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	opusSrcElement, err := pipeline.GetElementByName("opussrc")
	if err != nil {
		return err
	}
	opusSrc := app.SrcFromElement(opusSrcElement)

	opusSrc.SetCallbacks(&app.SourceCallbacks{
		NeedDataFunc: func(appsrc *app.Source, length uint) {
			buf := make([]byte, 1400)
			i, _, _ := sink.track.Read(buf)
			if i == 0 {
				return
			}
			buffer := gst.NewBufferWithSize(int64(i))
			buffer.Map(gst.MapWrite).WriteData(buf[:i])
			buffer.Unmap()
			appsrc.PushBuffer(buffer)
		},
	})

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)
	sink.pipeline = pipeline
	sink.loop = loop

	go sink.loop.Run()
	print("started")
	return nil
}

func (sink *RemoteAudioSink) Close() {
	if sink.loop == nil {
		return
	}
	sink.pipeline.BlockSetState(gst.StateNull)
	sink.loop.Quit()
}

func (sink *RemoteAudioSink) Track() *webrtc.TrackRemote {
	return sink.track
}

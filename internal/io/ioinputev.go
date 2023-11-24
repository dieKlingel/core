package io

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

type IOInputDevice struct {
	streams  map[string]*Stream
	pipeline *gst.Pipeline
	loop     *glib.MainLoop
	sink     *app.Sink
	name     string
}

func NewIOInputDevice(pipeline string) (*IOInputDevice, error) {
	pipe, err := gst.NewPipelineFromString(pipeline)
	if err != nil {
		return nil, err
	}

	if _, err := pipe.GetElementByName("rawsink"); err != nil {
		return nil, err
	}

	device := &IOInputDevice{
		streams:  make(map[string]*Stream),
		pipeline: pipe,
		name:     uuid.New().String(),
	}

	return device, nil
}

func (dev *IOInputDevice) SetName(name string) {
	dev.name = name
}

func (dev *IOInputDevice) GetName() string {
	return dev.name
}

func (dev *IOInputDevice) AddStream(stream *Stream) {
	if _, exists := dev.streams[stream.ID()]; exists {

		return
	}
	if stream.camera != nil {
		panic("a stream can only be added to one iodevice")
	}

	dev.streams[stream.ID()] = stream
	stream.camera = dev
	dev.openSafe()
}

func (dev *IOInputDevice) RemoveStream(stream *Stream) {
	delete(dev.streams, stream.ID())
	stream.Finished <- true
	stream.camera = nil
	dev.closeSafe()
}

func (dev *IOInputDevice) openSafe() {
	if dev.loop != nil && dev.loop.IsRunning() {
		return
	}
	log.Printf("open IOInputDevice %s", dev.name)

	dev.loop = glib.NewMainLoop(glib.MainContextDefault(), false)

	dev.pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
			dev.pipeline.BlockSetState(gst.StateNull)
			dev.loop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			dev.loop.Quit()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			// fmt.Println(msg)
		}
		return true
	})

	appsink, err := dev.pipeline.GetElementByName("rawsink")
	if err != nil {
		panic(err.Error())
	}
	dev.sink = app.SinkFromElement(appsink)

	dev.sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := dev.sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}

			for _, stream := range dev.streams {
				if !stream.IsReady() {
					continue
				}
				stream.src.PushSample(sample)
			}

			return gst.FlowOK
		},
	})

	dev.pipeline.SetState(gst.StatePlaying)
	go dev.loop.Run()
}

func (dev *IOInputDevice) closeSafe() {
	if len(dev.streams) != 0 {
		return
	}

	if dev.loop == nil {
		return
	}

	log.Printf("close IOInputDevice %s", dev.name)

	if dev.pipeline != nil {
		dev.pipeline.BlockSetState(gst.StateNull)
	}
	if dev.loop != nil {
		dev.loop.Quit()
	}
}

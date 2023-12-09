package camera

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type CameraCodec string

const (
	X264CameraCodec  CameraCodec = "appsrc name=src ! videoconvert !  x264enc tune=zerolatency bitrate=500 speed-preset=superfast ! appsink sync=false name=sink"
	MJPEGCameraCodec CameraCodec = "appsrc name=src ! videoconvert ! jpegenc ! appsink sync=false name=sink"
)

func (codec CameraCodec) ToPipelineElements() (*app.Source, *gst.Pipeline, *app.Sink) {
	pipeline, err := gst.NewPipelineFromString(string(codec))
	if err != nil {
		panic(err)
	}
	source, err := pipeline.GetElementByName("src")
	if err != nil {
		panic(err)
	}
	sink, err := pipeline.GetElementByName("sink")
	if err != nil {
		panic(err)
	}

	return app.SrcFromElement(source), pipeline, app.SinkFromElement(sink)
}

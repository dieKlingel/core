package audio

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type AudioCodec string
type AudioPlayerCodec string

const (
	OpusEncodeCodec       AudioCodec       = "appsrc name=src ! audioconvert ! opusenc ! appsink sync=false name=sink"
	OpusDecodePlayerCodec AudioPlayerCodec = "application/x-rtp, payload=127, encoding-name=OPUS ! rtpopusdepay ! decodebin"
)

func (codec AudioCodec) ToPipelineElements() (*app.Source, *gst.Pipeline, *app.Sink) {
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

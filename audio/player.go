package audio

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type Player struct {
	track    Track
	codec    AudioPlayerCodec
	pipeline *gst.Pipeline
	source   *app.Source
}

func NewPlayer(codec AudioPlayerCodec) (*Player, error) {
	player := &Player{
		codec: codec,
	}

	pipeline, err := gst.NewPipelineFromString(fmt.Sprintf("appsrc format=time do-timestamp=true ! %s ! autoaudiosink", string(codec)))
	if err != nil {
		return nil, err
	}
	player.pipeline = pipeline
	sources, err := pipeline.GetSourceElements()
	if err != nil {
		return nil, err
	}
	if len(sources) != 1 {
		panic(fmt.Sprintf("expected 1 source, got %d", len(sources)))
	}
	player.source = app.SrcFromElement(sources[0])

	return player, nil
}

func (player *Player) Play(track Track) {
	if player.track != nil {
		panic("player already playing")
	}
	player.track = track

	player.source.SetCallbacks(&app.SourceCallbacks{
		NeedDataFunc: func(appsrc *app.Source, length uint) {
			buf := make([]byte, 1400)
			i, _, _ := track.Read(buf)

			if i == 0 {
				return
			}

			buffer := gst.NewBufferWithSize(int64(i))
			buffer.Map(gst.MapWrite).WriteData(buf[:i])
			buffer.Unmap()
			appsrc.PushBuffer(buffer)
		},
	})

	player.pipeline.SetState(gst.StatePlaying)
}

func (player *Player) Stop() {
	if player.track == nil {
		return
	}
	player.pipeline.SetState(gst.StateNull)
	player.track = nil
}

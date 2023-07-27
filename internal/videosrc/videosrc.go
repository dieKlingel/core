package videosrc

import (
	"image"

	"github.com/google/uuid"
	"github.com/pion/mediadevices"
)

var refs = 0

func NewSharedVideoTrack(selector *mediadevices.CodecSelector) (mediadevices.Track, error) {
	camera.Open()

	src := NewSharedVideoSource()
	track := mediadevices.NewVideoTrack(
		&src,
		selector,
	)

	return track, nil
}

type SharedVideoSource struct {
	id string
}

func NewSharedVideoSource() SharedVideoSource {
	refs++
	return SharedVideoSource{
		id: uuid.New().String(),
	}
}

func (src *SharedVideoSource) Read() (image.Image, func(), error) {
	mat := camera.Read()

	relfunc := func() {
		mat.Close()
		// TODO release img
	}

	img, err := mat.ToImage()
	return img, relfunc, err
}

func (src *SharedVideoSource) ID() string {
	return src.id
}

func (src *SharedVideoSource) Close() error {
	refs--
	print("refs ")
	print(refs)
	if refs == 0 {
		camera.Close()
	}
	return nil
}

package io

type CameraCodec string

const (
	X264CameraCodec  CameraCodec = ""
	RawCameraCodec   CameraCodec = ""
	MJPEGCameraCodec CameraCodec = "appsrc name=src ! videoconvert ! jpegenc ! appsink sync=false name=sink"
)

type Camera interface {
	NewStream(codec CameraCodec) (*Stream, error)
	ReleaseStream(stream *Stream)
	HasStream() bool
}

type camera struct {
	iodev *IOInputDevice
}

func NewCamera(pipeline string) (Camera, error) {
	dev, err := NewIOInputDevice(pipeline)
	if err != nil {
		return nil, err
	}
	dev.SetName("io-camera")

	return &camera{
		iodev: dev,
	}, nil
}

func (cam *camera) NewStream(codec CameraCodec) (*Stream, error) {
	stream, err := NewStream(string(codec))
	if err != nil {
		return nil, err
	}

	cam.iodev.AddStream(stream)
	return stream, nil
}

func (cam *camera) ReleaseStream(stream *Stream) {
	cam.iodev.RemoveStream(stream)
}

func (cam *camera) HasStream() bool {
	return len(cam.iodev.streams) != 0
}

package io

type Microphone interface {
	NewStream(codec AudioCodec) (*Stream, error)
	ReleaseStream(stream *Stream)
	HasStream() bool
}

type microphone struct {
	iodev *IOInputDevice
}

func NewMicrophone(pipeline string) (Microphone, error) {
	dev, err := NewIOInputDevice(pipeline)
	if err != nil {
		return nil, err
	}
	dev.SetName("io-camera")

	return &microphone{
		iodev: dev,
	}, nil
}

func (mic *microphone) NewStream(codec AudioCodec) (*Stream, error) {
	stream, err := NewStream(string(codec))
	if err != nil {
		return nil, err
	}

	mic.iodev.AddStream(stream)
	return stream, nil
}

func (mic *microphone) ReleaseStream(stream *Stream) {
	mic.iodev.RemoveStream(stream)
}

func (mic *microphone) HasStream() bool {
	return len(mic.iodev.streams) != 0
}

package io

type AudioCodec string

const (
	OpusAudioCodec AudioCodec = "appsrc name=src ! audioconvert ! opusenc ! appsink sync=false name=sink"
)

type CameraCodec string

const (
	X264CameraCodec  CameraCodec = "appsrc name=src ! videoconvert !  x264enc tune=zerolatency bitrate=500 speed-preset=superfast ! appsink sync=false name=sink"
	MJPEGCameraCodec CameraCodec = "appsrc name=src ! videoconvert ! jpegenc ! appsink sync=false name=sink"
)

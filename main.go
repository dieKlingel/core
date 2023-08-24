package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dieklingel/core/internal/video"
)

var config *Config
var camera *video.Camera

func main() {
	wd := os.Getenv("DIEKLINGEL_HOME")

	if len(strings.TrimSpace(wd)) == 0 {
		log.Printf("the environment variable DIEKLINGEL_HOME is not set")
	} else if err := syscall.Chdir(wd); err != nil {
		log.Printf("error while switching to workdir: %s", err.Error())
		os.Exit(1)
	}

	dir, _ := syscall.Getwd()
	log.Printf("Running in working directory: %s", dir)

	conf, err := NewConfigFromCurrentDirectory()
	if err != nil {
		log.Printf("cannot read config fiel: %s", err.Error())
		os.Exit(1)
	}
	config = conf

	uri, err := url.Parse(config.Mqtt.Uri)
	if err != nil {
		log.Printf("cannot parse mqtt uri: %s", err.Error())
		os.Exit(1)
	}

	camera, err = video.NewCamera(config.Media.VideoSrc)
	if err != nil {
		log.Printf(`cannot create the camera from video-src: %s.
	A possible cause could be the upgrade to version 0.3.0 or higher.
	Since version 0.3.0 we no longer use 'h264sink' in our video-src pipeline.
	Instead we use 'rawsink' which should emit a raw video stream,
	which we will convert internal. In order to fix this, build your video-src pipeline like:
	...
	  media:
  	    video-src: autovideosrc ! video/x-raw, framerate=30/1, width=1280, height=720 ! appsink name=rawsink
	...`, err.Error(),
		)
	}

	RunApi(
		*uri,
		config.Mqtt.Username,
		config.Mqtt.Password,
	)
	RunProxy(8081)

	// Wait for interruption to exit
	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
}

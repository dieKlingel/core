package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dieklingel/core/endpoint"
	"github.com/dieklingel/core/internal/io"
	"github.com/dieklingel/core/service"
	"github.com/dieklingel/core/transport"
	"github.com/spf13/viper"
)

var config *Config
var camera *io.IOInputDevice
var microphone *io.IOInputDevice

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

	_, err = url.Parse(config.Mqtt.Uri)
	if err != nil {
		log.Printf("cannot parse mqtt uri: %s", err.Error())
		os.Exit(1)
	}

	camera, err = io.NewIOInputDevice(config.Media.VideoSrc)
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
	} else {
		camera.SetName("Camera")
	}

	microphone, err = io.NewIOInputDevice(config.Media.AudioSrc)
	if err != nil {
		log.Printf(`cannot create the microphone from audio-src: %s.
	A possible cause could be the upgrade to version 0.3.0 or higher.
	Since version 0.3.0 we no longer use 'opussink' in our auidio-src pipeline.
	Instead we use 'rawsink' which should emit a raw audio stream,
	which we will convert internal. In order to fix this, build your audio-src pipeline like:
	...
	  media:
  	    audio-src: autoaudiosrc ! audio/x-raw, format=S16LE, layout=interleaved, rate=48000, channels=1 ! appsink name=rawsink ! appsink name=rawsink
	...`, err.Error(),
		)
	} else {
		microphone.SetName("Microphone")
	}

	viper.SetConfigName("core")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$DIEKLINGEL_HOME")
	viper.AddConfigPath("/etc/dieklingel")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("the config file coulf not be found")
			os.Exit(1)
		} else {
			log.Fatal("the config file could not be parsed")
			os.Exit(2)
		}
	}

	viper.SetDefault("http.port", "8080")
	viper.SetDefault("mqtt.uri", "mqtts://server.dieklingel.com:8883/dieklingel/mayer/kai/")

	system := endpoint.NewSystemEndpoint(service.NewSystemService())
	action := endpoint.NewActionEndpoint(service.NewActionService())

	//action.Add("test", "echo H")

	transport.NewHttpTransport(8080, system, action).Run()

	// Wait for interruption to exit
	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	// TODO: cleanup
}

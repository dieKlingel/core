package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Environment struct {
	Actions []Action `yaml:"actions"`
	Mqtt    Mqtt     `yaml:"mqtt"`
	Rtc     Rtc      `yaml:"rtc"`
	Media   Media    `yaml:"media"`
}

func New() *Environment {
	file, err := os.Open("core.yaml")
	if err != nil {
		panic(err.Error())
	}

	env := &Environment{
		Actions: make([]Action, 0),
		Mqtt: Mqtt{
			Uri:      "mqtt://localhost:1883",
			Username: "guest",
			Password: "guest",
		},
		Rtc: Rtc{
			IceServers: []IceServer{
				{
					Urls: "stun:stun.l.google.com:19302",
				},
			},
		},
		Media: Media{
			Camera: Camera{
				Src: "videotestsrc ! video/x-raw, framerate=30/1, width=1280, height=720 ! appsink",
			},
			Audio: Audio{
				Src:  "audiotestsrc ! audio/x-raw, format=S16LE, layout=interleaved, rate=48000, channels=1 ! appsink",
				Sink: "autoaudiosink",
			},
		},
	}

	if err := yaml.NewDecoder(file).Decode(&env); err != nil {
		panic(err)
	}

	return env
}

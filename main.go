package main

import (
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/dieklingel/core/config"
	"go.uber.org/fx"
)

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

	app := fx.New(
		fx.Provide(
			config.New,
			NewFxCamera,
			NewFxAudioInput,
			NewActionService,
			NewFxHttpService,
			NewWebRTCService,
			NewMqttService,
		),
		fx.Invoke(
			func(h *HttpService, m *MqttService) {
				h.Run()
				m.Run()
			},
		),
	)
	app.Run()
}

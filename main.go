package main

import (
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/tinyzimmer/go-gst/gst"
	"go.uber.org/fx"
)

func main() {
	wd := os.Getenv("DIEKLINGEL_HOME")
	gst.Init(nil)

	if len(strings.TrimSpace(wd)) == 0 {
		log.Printf("the environment variable DIEKLINGEL_HOME is not set")
	} else if err := syscall.Chdir(wd); err != nil {
		log.Printf("error while switching to workdir: %s", err.Error())
		os.Exit(1)
	}

	dir, _ := syscall.Getwd()
	log.Printf("Running in working directory: %s", dir)

	fx.New(
		fx.Provide(
			NewFxStorageService,
			NewFxCameraService,
			NewFxActionService,
			NewFxHttpService,
			NewFxWebRTCService,
			NewFxMqttService,
		),
		fx.Invoke(func(h *HttpService, m *MqttService) {
			h.Run()
			m.Run()
		}),
	).Run()
}

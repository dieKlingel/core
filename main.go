package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tinyzimmer/go-gst/gst"
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

	storagesrv := NewStorageService("core.yaml")
	camerasrv := NewCameraService(storagesrv)
	actionsrv := NewActionService(storagesrv)
	httpsrv := NewHttpService(8080, storagesrv, camerasrv)
	webrtcsrv := NewWebRTCService(camerasrv)
	mqttsrv := NewMqttService(storagesrv, actionsrv, webrtcsrv)

	httpsrv.Run()
	mqttsrv.Run()

	// Wait for interruption to exit
	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	// TODO: cleanup
}

package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	config, err := NewConfigFromCurrentDirectory()
	if err != nil {
		log.Printf("cannot read config fiel: %s", err.Error())
		os.Exit(1)
	}

	uri, err := url.Parse(config.Mqtt.Uri)
	if err != nil {
		log.Printf("cannot parse mqtt uri: %s", err.Error())
		os.Exit(1)
	}

	RunApi(
		*uri,
		config.Mqtt.Password,
		config.Mqtt.Username,
	)
	RunProxy(8081)

	// Wait for interruption to exit
	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
}

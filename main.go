package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dieklingel/core/actionsrv"
	"github.com/dieklingel/core/camerasrv"
	"github.com/dieklingel/core/devicesrv"
	"github.com/dieklingel/core/httpsrv"
	"github.com/dieklingel/core/internal/io"
	"github.com/dieklingel/core/mqttsrv"
	"github.com/dieklingel/core/signsrv"
	"github.com/dieklingel/core/usersrv"
	"github.com/dieklingel/core/webrtcsrv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	db, err := gorm.Open(sqlite.Open("core.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	actionsrv := actionsrv.NewService(db)
	devicesrv := devicesrv.NewService(db)
	signsrv := signsrv.NewService(db)
	usersrv := usersrv.NewService(db)
	camerasrv := camerasrv.NewService(db)
	webrtcsrv := webrtcsrv.NewService(camerasrv)
	mqttsrv := mqttsrv.NewService(db, devicesrv, actionsrv, webrtcsrv)
	httpsrv := httpsrv.NewService(8080, actionsrv, devicesrv, signsrv, usersrv, camerasrv, mqttsrv)

	httpsrv.Run()

	// Wait for interruption to exit
	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	// TODO: cleanup
}

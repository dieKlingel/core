package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/io"
	"github.com/gorilla/mux"
)

type HttpService struct {
	port           int
	cameraService  *CameraService
	storageService core.StorageService

	server *http.Server
}

func NewHttpService(port int, storageService core.StorageService, cameraService *CameraService) *HttpService {
	return &HttpService{
		port:           port,
		cameraService:  cameraService,
		storageService: storageService,
	}
}

func (httpService *HttpService) Run() error {
	router := mux.NewRouter()

	router.Methods("GET").Path("/camera/snapshot").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stream := httpService.cameraService.NewCameraStream(io.MJPEGCameraCodec)
		if stream == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer httpService.cameraService.ReleaseCameraStream(stream)

		select {
		case frame := <-stream.Frame:
			w.Header().Add("Content-Type", "image/jpeg")
			img := frame.GetBuffer().Bytes()
			n, err := w.Write(img)
			if err != nil || n != len(img) {
				return
			}
		case <-r.Context().Done():
			return
		case <-time.After(5 * time.Second):
			log.Println("the http mjpeg stream was closed by timeout of 5 seconds, cause no frame could be received but the connection was still open")
			return
		}

	})

	router.Methods("GET").Path("/camera/stream").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stream := httpService.cameraService.NewCameraStream(io.MJPEGCameraCodec)
		if stream == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer httpService.cameraService.ReleaseCameraStream(stream)

		w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		boundary := "\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"

		for {
			select {
			case frame := <-stream.Frame:
				img := frame.GetBuffer().Bytes()

				n, err := w.Write([]byte(boundary))
				if err != nil || n != len(boundary) {
					return
				}

				n, err = w.Write(img)
				if err != nil || n != len(img) {
					return
				}

				n, err = w.Write([]byte("\r\n"))
				if err != nil || n != 2 {
					return
				}
			case <-r.Context().Done():
				return
			case <-time.After(5 * time.Second):
				log.Println("the http mjpeg stream was closed by timeout of 5 seconds, cause no frame could be received but the connection was still open")
				return
			}
		}
	})
	httpService.server = &http.Server{
		Handler:     router,
		Addr:        fmt.Sprintf(":%d", httpService.port),
		ReadTimeout: 15 * time.Second,
	}

	go httpService.server.ListenAndServe()
	return nil
}

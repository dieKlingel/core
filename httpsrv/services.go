package httpsrv

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	camio "github.com/dieklingel/core/internal/io"
	"github.com/gorilla/mux"
)

func buildServiceRoutes(service *HttpService, router *mux.Router) {
	router.Methods("GET").Path("/camera").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			CameraPipeline string
		}

		json.NewEncoder(w).Encode(Response{
			CameraPipeline: service.CameraService.CameraPipeline(),
		})
	})

	router.Methods("GET").Path("/camera/stream").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stream := service.CameraService.NewCameraStream(camio.MJPEGCameraCodec)
		if stream == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer service.CameraService.ReleaseCameraStream(stream)

		w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		boundary := "\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"

		for {
			select {
			case frame := <-stream.Frame:
				img := frame.GetBuffer().Bytes()

				n, err := io.WriteString(w, boundary)
				if err != nil || n != len(boundary) {
					return
				}

				n, err = w.Write(img)
				if err != nil || n != len(img) {
					return
				}

				n, err = io.WriteString(w, "\r\n")
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
}

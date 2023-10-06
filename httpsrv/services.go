package httpsrv

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dieklingel/core/internal/core"
	camio "github.com/dieklingel/core/internal/io"
	"github.com/dieklingel/core/internal/slice"
	"github.com/gorilla/mux"
)

func buildServiceRoutes(service *HttpService, router *mux.Router) {
	type MqttConnectionResponse struct {
		Id           uint64
		Url          string
		Username     string
		IsConnected  bool
		ErrorMessage string
	}

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

	router.Methods("GET").Path("/mqtt/connections").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		connections := slice.Map(service.MqttService.Connections(), func(c core.MqttConnection) MqttConnectionResponse {
			return MqttConnectionResponse{
				Id:           c.Id,
				Url:          c.Url,
				Username:     c.Username,
				IsConnected:  c.Client.IsConnected(),
				ErrorMessage: c.Client.ErrorMessage(),
			}
		})

		json.NewEncoder(w).Encode(connections)
	})

	router.Methods("PUT").Path("/mqtt/connections").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Url      string
			Username string
			Password string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		connection := core.MqttConnection{
			Url:      req.Url,
			Username: req.Username,
			Password: req.Password,
		}
		service.MqttService.SaveConnection(&connection)
		json.NewEncoder(w).Encode(MqttConnectionResponse{
			Id:           connection.Id,
			Url:          connection.Url,
			Username:     connection.Username,
			IsConnected:  connection.Client.IsConnected(),
			ErrorMessage: connection.Client.ErrorMessage(),
		})
	})

	router.Methods("PATCH").Path("/mqtt/connections/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Url      string
			Username string
			Password string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		connection := service.MqttService.GetConnectionById(int(id))
		if connection == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if len(req.Url) != 0 {
			connection.Url = req.Url
		}
		if len(req.Username) != 0 {
			connection.Username = req.Username
		}
		if len(req.Password) != 0 {
			connection.Password = req.Password
		}
		service.MqttService.SaveConnection(connection)
		res := MqttConnectionResponse{
			Id:           connection.Id,
			Url:          connection.Url,
			Username:     connection.Username,
			IsConnected:  connection.Client.IsConnected(),
			ErrorMessage: connection.Client.ErrorMessage(),
		}
		json.NewEncoder(w).Encode(res)
	})

	router.Methods("DELETE").Path("/mqtt/connections/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		if connection := service.MqttService.GetConnectionById(int(id)); connection != nil {
			service.MqttService.RemoveConnection(connection)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

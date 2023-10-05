package httpsrv

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dieklingel/core/internal/core"
	"github.com/gorilla/mux"
)

func buildDeviceRoutes(service *HttpService, router *mux.Router) {
	router.Methods("GET").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		devices := service.DeviceService.Devices()

		payload, err := json.Marshal(devices)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	})

	router.Methods("PUT").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Token string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		device := &core.Device{}
		service.DeviceService.SaveDevice(device)
		json.NewEncoder(w).Encode(device)
	})

	router.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		device := service.DeviceService.GetDeviceById(int(id))
		if device == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(device)
	})

	router.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		device := service.DeviceService.GetDeviceById(int(id))
		if device != nil {
			service.DeviceService.RemoveDevice(device)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

package httpsrv

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dieklingel/core/internal/core"
	"github.com/gorilla/mux"
)

func buildSignRoutes(service *HttpService, router *mux.Router) {
	router.Methods("GET").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signs := service.SignService.Signs()

		payload, err := json.Marshal(signs)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	})

	router.Methods("PUT").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Script string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		sign := &core.Sign{
			Script: req.Script,
		}
		service.SignService.SaveSign(sign)
		json.NewEncoder(w).Encode(sign)
	})

	router.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		sign := service.SignService.GetSignById(int(id))
		if sign == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(sign)
	})

	router.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		sign := service.SignService.GetSignById(int(id))
		if sign != nil {
			service.SignService.RemoveSign(sign)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

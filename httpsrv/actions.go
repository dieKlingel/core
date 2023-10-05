package httpsrv

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dieklingel/core/internal/core"
	"github.com/gorilla/mux"
)

func buildActionRoutes(service *HttpService, router *mux.Router) {
	router.Methods("GET").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actions := service.ActionService.Actions()

		payload, err := json.Marshal(actions)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	})

	router.Methods("PUT").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Trigger     string
			Script      string
			Environment core.ActionExecutionEnvironment
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result := service.ActionService.SaveAction(core.Action{
			Trigger:     req.Trigger,
			Script:      req.Script,
			Environment: req.Environment,
		})
		json.NewEncoder(w).Encode(result)
	})

	router.Methods("POST").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Pattern     string
			Environment map[string]string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		results := service.ActionService.Execute(req.Pattern, req.Environment)
		json.NewEncoder(w).Encode(results)
	})

	router.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		action := service.ActionService.GetActionById(int(id))
		if action == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(action)
	})

	router.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		action := service.ActionService.GetActionById(int(id))
		if action != nil {
			service.ActionService.RemoveAction(*action)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.Methods("POST").Path("/{id:[0-9]+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		action := service.ActionService.GetActionById(int(id))
		if action == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		type Request struct {
			Environment map[string]string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result := action.Execute(req.Environment)
		json.NewEncoder(w).Encode(result)
	})
}

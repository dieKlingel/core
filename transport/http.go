package transport

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dieklingel/core/internal/api"
	"github.com/dieklingel/core/transport/dashboard"
	"github.com/gorilla/mux"
)

type SystemEndpoint interface {
	Version() string
}

type ActionEndpoint interface {
	List() []api.Action
	Execute(pattern string, environment map[string]string) []api.ActionExecutionResult
	GetById(id string) api.Action
	Add(trigger string, script string) api.Action
}

type HttpTransport struct {
	port   int
	system SystemEndpoint
	action ActionEndpoint

	server *http.Server
}

func NewHttpTransport(port int, system SystemEndpoint, action ActionEndpoint) *HttpTransport {
	return &HttpTransport{
		port:   port,
		system: system,
		action: action,
	}
}

func (transport *HttpTransport) Port() int {
	return transport.port
}

func (transport *HttpTransport) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
		version := transport.system.Version()
		w.Write([]byte(version))
	})

	router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFS(dashboard.Files(), "html/index.html"))
		templ.Execute(w, nil)
	})

	router.HandleFunc("/dashboard/actions", func(w http.ResponseWriter, r *http.Request) {
		actions := transport.action.List()

		templ := template.Must(template.ParseFS(dashboard.Files(), "html/actions.html"))
		templ.Execute(w, actions)
	}).Methods("GET")

	router.HandleFunc("/dashboard/actions", func(w http.ResponseWriter, r *http.Request) {
		trigger := r.FormValue("trigger")
		script := r.FormValue("script")

		if len(trigger) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a value for the key 'trigger' was not provied"))
			return
		}
		if len(script) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a value for the key 'script' was not provied"))
			return
		}

		action := transport.action.Add(trigger, script)

		templ := template.Must(template.ParseFS(dashboard.Files(), "html/action.html"))
		templ.Execute(w, action)
	}).Methods("POST")

	router.HandleFunc("/dashboard/actions/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		action := transport.action.GetById(vars["id"])

		if action == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		templ := template.Must(template.ParseFS(dashboard.Files(), "html/action.html"))
		templ.Execute(w, action)
	})

	transport.server = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", transport.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go transport.server.ListenAndServe()
	return nil
}

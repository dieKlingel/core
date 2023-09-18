package transport

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dieklingel/core/transport/dashboard"
	"github.com/gorilla/mux"
)

type SystemEndpoint interface {
	Version() string
}

type HttpTransport struct {
	port   int
	system SystemEndpoint
	server *http.Server
}

func NewHttpTransport(port int, system SystemEndpoint) *HttpTransport {
	return &HttpTransport{
		port:   port,
		system: system,
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

	transport.server = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", transport.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go transport.server.ListenAndServe()
	return nil
}

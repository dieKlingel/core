package transport

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	hw "github.com/dieklingel/core/internal/io"
	"github.com/dieklingel/core/transport/dashboard"
	"github.com/gorilla/mux"
)

type HttpTransport struct {
	port   int
	system SystemEndpoint
	action ActionEndpoint
	sign   SignEndpoint
	camera hw.Camera

	server *http.Server
}

func NewHttpTransport(port int, system SystemEndpoint, action ActionEndpoint, sign SignEndpoint, camera hw.Camera) *HttpTransport {
	return &HttpTransport{
		port:   port,
		system: system,
		action: action,
		sign:   sign,
		camera: camera,
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
	}).Methods("GET")

	router.HandleFunc("/dashboard/actions/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		action := transport.action.GetById(vars["id"])

		if action == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		transport.action.Delete(action)
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	router.HandleFunc("/dashboard/signs", func(w http.ResponseWriter, r *http.Request) {
		signs := transport.sign.List()

		templ := template.Must(template.ParseFS(dashboard.Files(), "html/signs.html"))
		templ.Execute(w, signs)
	}).Methods("GET")

	router.HandleFunc("/dashboard/signs", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		script := r.FormValue("script")

		if len(name) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a value for the key 'name' was not provied"))
			return
		}
		if len(script) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a value for the key 'script' was not provied"))
			return
		}

		sign := transport.sign.Add(name, script)

		templ := template.Must(template.ParseFS(dashboard.Files(), "html/sign.html"))
		templ.Execute(w, sign)
	}).Methods("POST")

	router.HandleFunc("/dashboard/signs/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sign := transport.sign.GetById(vars["id"])

		if sign == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		templ := template.Must(template.ParseFS(dashboard.Files(), "html/sign.html"))
		templ.Execute(w, sign)
	}).Methods("GET")

	router.HandleFunc("/dashboard/signs/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sign := transport.sign.GetById(vars["id"])

		if sign == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		transport.sign.Delete(sign)
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	router.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		stream, err := transport.camera.NewStream(hw.MJPEGCameraCodec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer transport.camera.ReleaseStream(stream)

		w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		boundary := "\r\n--frame\r\nContent-Type: image/jpeg\r\n\r\n"

		for {
			select {
			case frame := <-stream.Frame:
				img := frame.GetBuffer().Bytes()

				if err != nil {
					return
				}

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
	}).Methods("GET")

	router.HandleFunc("/snapshot", func(w http.ResponseWriter, r *http.Request) {
		stream, err := transport.camera.NewStream(hw.MJPEGCameraCodec)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer transport.camera.ReleaseStream(stream)

		select {
		case frame := <-stream.Frame:
			img := frame.GetBuffer().Bytes()
			if err != nil {
				return
			}

			n, err := w.Write(img)
			if err != nil || n != len(img) {
				return
			}
		case <-time.After(5 * time.Second):
			log.Println("the http jpeg snapshot was closed by timeout of 5 seconds, cause no frame could be received but the connection was still open")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	transport.server = &http.Server{
		Handler:     router,
		Addr:        fmt.Sprintf(":%d", transport.port),
		ReadTimeout: 15 * time.Second,
	}

	go transport.server.ListenAndServe()
	return nil
}

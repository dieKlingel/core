package httpsrv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dieklingel/core/internal/core"
	"github.com/gorilla/mux"
)

type HttpService struct {
	Port          int
	ActionService core.ActionService
	DeviceService core.DeviceService
	SignService   core.SignService
	UserService   core.UserService
	CameraService core.CameraService

	server *http.Server
}

func NewService(port int, actionsrv core.ActionService, devicesrv core.DeviceService, signsrv core.SignService, usersrv core.UserService, camerasrv core.CameraService) core.HttpService {
	return &HttpService{
		Port:          port,
		ActionService: actionsrv,
		DeviceService: devicesrv,
		SignService:   signsrv,
		UserService:   usersrv,
		CameraService: camerasrv,
	}
}

func (transport *HttpService) Run() error {
	router := mux.NewRouter()

	//router.NewRoute().Handler(createActionsRouter(transport))
	buildActionRoutes(transport, router.PathPrefix("/actions").Subrouter())
	buildDeviceRoutes(transport, router.PathPrefix("/devices").Subrouter())
	buildSignRoutes(transport, router.PathPrefix("/signs").Subrouter())
	buildUserRoutes(transport, router.PathPrefix("/users").Subrouter())
	buildServiceRoutes(transport, router.PathPrefix("/services").Subrouter())

	transport.server = &http.Server{
		Handler:     router,
		Addr:        fmt.Sprintf(":%d", transport.Port),
		ReadTimeout: 15 * time.Second,
	}

	go transport.server.ListenAndServe()
	return nil
}

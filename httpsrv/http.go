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

	server *http.Server
}

func NewService(port int, actionsrv core.ActionService, devicesrv core.DeviceService) core.HttpService {
	return &HttpService{
		Port:          port,
		ActionService: actionsrv,
		DeviceService: devicesrv,
	}
}

func (transport *HttpService) Run() error {
	router := mux.NewRouter()

	//router.NewRoute().Handler(createActionsRouter(transport))
	buildActionRoutes(transport, router.PathPrefix("/actions").Subrouter())
	buildDeviceRoutes(transport, router.PathPrefix("/devices").Subrouter())

	transport.server = &http.Server{
		Handler:     router,
		Addr:        fmt.Sprintf(":%d", transport.Port),
		ReadTimeout: 15 * time.Second,
	}

	go transport.server.ListenAndServe()
	return nil
}

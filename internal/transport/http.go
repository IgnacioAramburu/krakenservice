package transport

import (
	"log"
	"ltp/internal/api"
	"ltp/internal/config"
	"ltp/internal/service"
	"net/http"
)

// NewHTTPRouter initializes the router using the services as dependencies to build the handlers.
func NewHTTPRouter(svc service.Service, conf config.Config, logger *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	srvHdlr := api.NewHandler(svc, conf, logger)

	mux.HandleFunc("/api/v1/ltp", srvHdlr.LastTradePriceHandler)

	return mux
}

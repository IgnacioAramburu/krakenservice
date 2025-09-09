package api

import (
	"log"
	"ltp/internal/_errors"
	"ltp/internal/config"
	"ltp/internal/service"
	"net/http"
	"strings"
)

type Handler struct {
	service service.Service
	conf    config.Config
	log     *log.Logger
}

func (h *Handler) LastTradePriceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()

	h.log.Printf("%s %s -> %s", r.Method, r.URL.Path, r.URL.RawQuery)

	pairInfo := ""
	pairQuery, ok := query["pair"]

	if ok && len(pairQuery) > 0 {
		pairInfo = strings.TrimSpace(pairQuery[0])
	}

	ltp, errInfo := h.service.GetLastTradedPrice(ctx, pairInfo)

	if errInfo != nil {
		switch errInfo.Type {
		case _errors.InvalidInputErrType:
			_ = RespondWithError(w, BadRequestErr.WithDescription(errInfo.Description))
		case _errors.GatewayErrType:
			_ = RespondWithError(w, GatewayErr)
		default:
			_ = RespondWithError(w, InternalServerErr)
			return
		}
	}

	_ = RespondWithData(w, http.StatusOK, ltp)

}

func NewHandler(svc service.Service, conf config.Config, logger *log.Logger) *Handler {
	return &Handler{
		service: svc,
		conf:    conf,
		log:     logger,
	}
}

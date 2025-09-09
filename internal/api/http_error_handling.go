package api

import (
	"fmt"
	"ltp/internal/_errors"
	"net/http"
)

var (
	InternalServerErr = APIError{
		Code: http.StatusInternalServerError,
		ErrorInfo: _errors.ErrorInformation{
			Type: _errors.InternalErrType,
		},
	}
	GatewayErr = APIError{
		Code: http.StatusBadGateway,
		ErrorInfo: _errors.ErrorInformation{
			Type:        _errors.GatewayErrType,
			Description: "error at communicating with trade info provider.",
		},
	}
	BadRequestErr = APIError{
		Code: http.StatusBadRequest,
		ErrorInfo: _errors.ErrorInformation{
			Type: _errors.InvalidInputErrType,
		},
	}
)

type APIError struct {
	Code      int `json:"code"`
	ErrorInfo _errors.ErrorInformation
}

func (e APIError) WithDescription(description string) APIError {
	return APIError{
		Code: e.Code,
		ErrorInfo: _errors.ErrorInformation{
			Type:        e.ErrorInfo.Type,
			Description: description,
		},
	}
}

func (e APIError) WithParams(params ...interface{}) APIError {
	return APIError{
		Code: e.Code,
		ErrorInfo: _errors.ErrorInformation{
			Type:        e.ErrorInfo.Type,
			Description: fmt.Sprintf(e.ErrorInfo.Description, params...),
		},
	}
}

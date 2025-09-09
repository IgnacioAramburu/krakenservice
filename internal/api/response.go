package api

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Result any `json:"result,omitempty"`
	Error  any `json:"error,omitempty"`
}

func newBaseResponseWithData(data any) BaseResponse {
	return BaseResponse{
		Result: data,
	}
}

func newBaseResponseWithError(err any) BaseResponse {
	return BaseResponse{
		Error: err,
	}
}

func RespondWithData(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, errObject APIError) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errObject.Code)
	return json.NewEncoder(w).Encode(newBaseResponseWithError(errObject.ErrorInfo))
}

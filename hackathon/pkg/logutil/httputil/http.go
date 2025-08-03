package httputil

import (
	"encoding/json"
	"net/http"
)

func ResponseInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(internalServerError())
}

func internalServerError() *Response {
	return &Response{
		Error: ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		},
	}
}

func ResponseMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(methodNotAllowed())
}

func methodNotAllowed() *Response {
	return &Response{
		Error: ErrorResponse{
			Code:    "METHOD_NOT_ALLOWED",
			Message: "method not allowed"},
	}
}

func ResponseBadRequest(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := &Response{
		Error: ErrorResponse{
			Code:    code,
			Message: message,
		},
	}
	json.NewEncoder(w).Encode(resp)
}

type Response struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

package web

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type response map[string]interface{}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

// OK sends successful JSON response
func OK(data interface{}, w http.ResponseWriter) {
	jsonResponse(data, http.StatusOK, w)
}

// Error sends error response
func Error(err error, code int, w http.ResponseWriter) {
	jsonResponse(errorResponse{Code: code, Message: err.Error()}, code, w)
}

func jsonResponse(data interface{}, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.Error(err)
		return
	}
}

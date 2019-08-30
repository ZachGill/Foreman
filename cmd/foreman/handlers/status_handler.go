package handlers

import (
	"net/http"
)

// Status provides a status for the server
type Status struct{}

// ServeHTTP always returns an OK response to indicate the server is running
func (handler *Status) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Server is running"))
}

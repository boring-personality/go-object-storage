package handlers

import (
	"fmt"
	"net/http"
)

type MonitorHandler struct{}

func NewMonitorHandler() *MonitorHandler {
	return &MonitorHandler{}
}

func (m *MonitorHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "OK")
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, jsonResponse{Message: "Method not allowed"})
	}
	resp := jsonResponse {
		Message: "Status OK",
	}
	writeJson(w, http.StatusOK, resp)
	fmt.Println(r.Method, "/health ", http.StatusOK)
}

func (m *MonitorHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, jsonResponse{Message: "Method not allowed"})
		return
	}
	http.ServeFile(w, r, "index.html")
}

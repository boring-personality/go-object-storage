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
	fmt.Println(r.Method, "/health")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

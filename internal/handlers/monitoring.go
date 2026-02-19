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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Println(r.Method, "/health")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

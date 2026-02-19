package main

import (
	"fmt"
	"net/http"
	"github.com/boring-personality/go-object-storage/internal/handlers"
)

const PORT = "8001"

func main() {
	// Creating a servermux
	mux := http.NewServeMux()

	monitorHandler := handlers.NewMonitorHandler()
	// add the endpoint
	mux.HandleFunc("GET /health", monitorHandler.HealthHandler)
	mux.HandleFunc("GET /", serveIndex)
	storageHandle := handlers.NewStorageHandler()
	mux.HandleFunc("POST /upload", storageHandle.UploadFile)

	// get file via id endpoint
	mux.HandleFunc("GET /download/{id}", storageHandle.DownloadFile)
	fmt.Printf("Starting the server at %s\n", PORT)
	port_string := fmt.Sprintf(":%s", PORT)
	err := http.ListenAndServe(port_string, mux)
	if err != nil {
		fmt.Printf("Error is starting the http server: %s\n", err)
	}
	defer storageHandle.Data.DB.Close()
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

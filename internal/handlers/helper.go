package handlers

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Message string 	`json:"message,omitempty"`
	ID		string	`json:"id,omitempty"`
}

func writeJson(w http.ResponseWriter, status int, data any, header ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(header) > 0 {
		for key, val := range(header[0]) {
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

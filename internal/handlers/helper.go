package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"encoding/json"
	"net/http"
	"path/filepath"
	"mime/multipart"
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

func storeFile(filename string, f multipart.File) (string, string, error) {
	token_string, _ := get_file_token(10)
	extension := filepath.Ext(filename)
	dst_path := filepath.Join("./data", token_string+extension)	// need to think of something to store the destination path
	dst, err :=	os.Create(dst_path)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, f)
	return token_string, dst_path, err
}

func get_file_token(token_length int) (string, error) {
	token := make([]byte, token_length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	token_string := hex.EncodeToString(token)
	return token_string, nil
}

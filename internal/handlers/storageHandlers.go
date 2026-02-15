package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// var dictionary = make(map[string]string)

type StorageHandler struct{}

func NewStorageHandler() *StorageHandler {
	return &StorageHandler{}
}

func (sh *StorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, "/upload")
	err := r.ParseMultipartForm(10<<20) // This value approximately comes out to be 10MB
	if err != nil {
		http.Error(w, "Some issue uploading the file", http.StatusBadRequest)
	}

	file, header, err := r.FormFile("File")
	if err != nil {
		http.Error(w, "Error fetching the file", http.StatusBadRequest)
	}

	fmt.Printf("Filename: %s, Size: %d\n", header.Filename, header.Size)
	defer file.Close()
	
	token_length := 10
	token := make([]byte, token_length)
	rand.Read(token)
	token_string := hex.EncodeToString(token)
	
	extension := filepath.Ext(header.Filename)
	dst_path := filepath.Join("./data", token_string+extension)
	dst, err :=	os.Create(dst_path)
	if err != nil {
		fmt.Println("Error in saving file", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error in writing to the file", err)
	}

	// let's keep the token to file mapping in non persistent memory for now
	// dictionary[token_string] = dst_path

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "The file is uploaded succesfully. Here is the ID: %s", token_string)
}

func (sh *StorageHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("id")
	fmt.Println(r.Method, "/download/", token)
	// fmt.Printf("The requested file location is %s\n", dictionary[token])

	src_path := filepath.Join("./data", token)
	matches, _ := filepath.Glob(src_path+".*")
	if len(matches) == 0 {
		http.Error(w, "Failed to locate the file", http.StatusBadRequest)
		return
	}
	dst, err := os.Open(matches[0])
	if err != nil {
		http.Error(w, "Failed to locate the file", http.StatusBadRequest)
		return
	}
	defer dst.Close()
	
	extension := filepath.Ext(matches[0])
	// this tells the browser to treat the data as sequence of 8 bit bytes
	// so that the browser does not try to render it
	w.Header().Set("Content-Type", "application/octet-stream")

	// this tell the Content-Disposition is of type attachment so download the file with filename given
	w.Header().Set("Content-Disposition", "attachment; filename="+token+extension)

	_, err = io.Copy(w, dst)
	if err != nil {
		fmt.Println("Failed to send the file to client", err)
		return
	}
}

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/boring-personality/go-object-storage/internal/db"
)

type StorageHandler struct{
	Data *db.Database
}

func NewStorageHandler() *StorageHandler {
	do, err := db.NewDatabase()
	if err != nil {
		fmt.Println(err)
	}
	err = do.Create()
	if err != nil {
		fmt.Println(err)
	}
	return &StorageHandler{Data: do,}
}

func (sh *StorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, jsonResponse{Message: "Method not allowed"})
	}

	fmt.Println(r.Method, "/upload")
	err := r.ParseMultipartForm(10<<20) // This value approximately comes out to be 10MB
	if err != nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Some issue uploading the file"})
		return
	}

	file, header, err := r.FormFile("File")
	if err != nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Error fetching the file"})
		return
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
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Error saving file"})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Error fetching the file"})
		return
	}

	var obj db.Object
	obj.Id = token_string
	obj.Disk_path = dst_path
	obj.Original_name = header.Filename
	obj.Size = header.Size

	err = sh.Data.Insert(obj)
	resp := jsonResponse {
		Message: 	"The file is updated succesfully",
		ID:			token_string,
	}
	writeJson(w, http.StatusOK, resp)
}

func (sh *StorageHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, jsonResponse{Message: "Method not allowed"})
	}

	token := r.PathValue("id")
	fmt.Println(r.Method, "/download/", token)
	// fmt.Printf("The requested file location is %s\n", dictionary[token])

	obj, err := sh.Data.Read(token)
	if err != nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Failed to get file from database"})
		return
	}
	if obj == nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "File not present in the database"})
		return
	}

	dst, err := os.Open(obj.Disk_path)
	if err != nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Failed to locate the file"})
		return
	}
	defer dst.Close()

	// this tells the browser to treat the data as sequence of 8 bit bytes
	// so that the browser does not try to render it
	w.Header().Set("Content-Type", "application/octet-stream")

	// this tell the Content-Disposition is of type attachment so download the file with filename given
	w.Header().Set("Content-Disposition", "attachment; filename="+obj.Original_name)

	_, err = io.Copy(w, dst)
	if err != nil {
		writeJson(w, http.StatusBadRequest, jsonResponse{Message: "Failed to send the file to client"})
		return
	}
}

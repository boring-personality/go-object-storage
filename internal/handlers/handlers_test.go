package handlers

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"mime/multipart"
	"errors"
	"os"

	"github.com/boring-personality/go-object-storage/internal/db"
)

type FakeDatabase struct {
	Err 			error
	InsertCalled	bool
	InsertObj		db.Object
	ReadCalled		bool
	ReadObj			db.Object
}

func (fd *FakeDatabase) Insert(obj db.Object) error {
	fd.InsertCalled = true
	fd.InsertObj = obj
	return fd.Err
}

func (fd *FakeDatabase) Read(id string) (*db.Object, error) {
	fd.ReadCalled = true
	return &fd.ReadObj, fd.Err
}

func TestUploadFile_OK(t *testing.T) {
	fakedata := &FakeDatabase {}

	sh := &StorageHandler {
		Data: fakedata,
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("File", "test.txt")
	if err != nil {
		t.Fatal("Failed to set the test file metadata", err)
	}

	part.Write([]byte("Testing the upload file"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	sh.UploadFile(rr, req)

	if rr.Code != http.StatusCreated {
		t.Log(rr.Body.String())
		t.Fatal("Status code mis-match")
	}

	if !fakedata.InsertCalled {
		t.Fatal("failed to call insertdata")
	}

	if fakedata.InsertObj.Original_name != "test.txt" {
		t.Fatal("Unexpected filename", fakedata.InsertObj.Original_name)
	}
	os.Remove(fakedata.InsertObj.Disk_path)
}

func TestUploadFile_DBError(t *testing.T) {
	fakedata := &FakeDatabase {
		Err: errors.New("Database down"),
	}

	sh := &StorageHandler {
		Data: fakedata,
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("File", "test.txt")
	if err != nil {
		t.Fatal("Failed to set the test file metadata", err)
	}

	part.Write([]byte("Testing databse error"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	sh.UploadFile(rr, req)
	// this is very not correct but I am new to this so,
	// I will do it like this for on will fix this afterwards
	os.Remove(fakedata.InsertObj.Disk_path)

	if rr.Code != http.StatusInternalServerError {
		t.Fatal("status code mistmatch, the database should be down")
	}
}

func TestDownloadFile_OK(t *testing.T) {
	temp, _ := os.CreateTemp("", "test*")
	defer os.Remove(temp.Name())

	content := []byte("Testing download")
	temp.Write(content)
	temp.Close()

	obj := &db.Object {
		Id: 			"123",
		Original_name:	"test.txt",
		Disk_path:		temp.Name(),
		Size:			int64(len(content)),
	}

	fakedata := &FakeDatabase {
		ReadObj: *obj,
	}

	sh := &StorageHandler {
		Data: fakedata,
	}

	req := httptest.NewRequest(http.MethodGet, "/download/123", nil)
	rr := httptest.NewRecorder()

	sh.DownloadFile(rr, req)

	// if rr.Code != http.StatusOK {
	// 	t.Fatal("Failed to fetch the data")
	// }
	//
	if !fakedata.ReadCalled {
		t.Fatal("failed to call readdata")
	}

	if rr.Body.String() != string(content) {
		t.Fatal("File content don't match")
	}
	cd := rr.Header().Get("Content-Disposition")
	if cd == "" {
		t.Fatal("Content Disposition is not set")
	}
}

func TestDownloadFile_InvalidKey(t *testing.T) {
	temp, _ := os.CreateTemp("", "test*")
	defer os.Remove(temp.Name())

	content := []byte("Testing download")
	temp.Write(content)
	temp.Close()

	obj := &db.Object {
		Id: 			"123",
		Original_name:	"test.txt",
		Disk_path:		temp.Name(),
		Size:			int64(len(content)),
	}

	fakedata := &FakeDatabase {
		ReadObj: 	*obj,
		Err:		errors.New("Key not present in the database"),
	}

	sh := &StorageHandler {
		Data: fakedata,
	}

	req := httptest.NewRequest(http.MethodGet, "/download/1234", nil)
	rr := httptest.NewRecorder()

	sh.DownloadFile(rr, req)
	if !fakedata.ReadCalled {
		t.Fatal("failed to call readdata")
	}

	if rr.Code != http.StatusBadRequest {
		t.Fatal("the key should not be present in the database")
	}
}

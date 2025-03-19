package api

import (
	"encoding/json"
	"github.com/k3nnee/dfs/backend/master_node/app/schema"
	"io"
	"net/http"
)

const maxFileSize = 10 << 20

var allowedFileTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	err := r.ParseMultipartForm(maxFileSize)

	if err != nil {
		http.Error(w, "Unable to parse file: "+err.Error(), http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Unable to fetch file: "+err.Error(), http.StatusMethodNotAllowed)
		return
	}

	defer file.Close()

	if fileHeader.Size > maxFileSize {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	buf := make([]byte, 512)
	_, err = file.Read(buf)

	if err != nil {
		http.Error(w, "Error reading first 512 bytes: "+err.Error(), http.StatusBadRequest)
		return
	}

	fileType := http.DetectContentType(buf)

	if !allowedFileTypes[fileType] {
		http.Error(w, "File type not allowed: "+fileType, http.StatusMethodNotAllowed)
		return
	}

	file.Seek(0, io.SeekStart)

	// TODO: Splice data
	//fileContent, err := io.ReadAll(file)
	//
	//if err != nil {
	//	http.Error(w, "Unable to read file: "+err.Error(), http.StatusInternalServerError)
	//	return
	//}

	res := schema.Response{
		Message: "File uploaded successfully!",
		Success: true,
	}

	encodedMsg, err := json.Marshal(res)

	if err != nil {
		http.Error(w, "Unable to parse message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedMsg)
}

func UploadData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handle(w, r)
		return
	default:
		http.Error(w, r.Method+": method not allowed", http.StatusMethodNotAllowed)
	}
}

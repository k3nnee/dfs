package api

import (
	"encoding/json"
	"fmt"
	"github.com/k3nnee/dfs/backend/worker_node/app/constant"
	"github.com/k3nnee/dfs/backend/worker_node/app/schema"
	"io"
	"net/http"
	"os"
)

func UploadData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var fileUpload schema.FileUpload

	err = json.Unmarshal(body, &fileUpload)

	if err != nil {
		http.Error(w, "Error parsing body: "+err.Error(), http.StatusBadRequest)
		return
	}

	filePath := fileUpload.FileName + "." + constant.FileTypes[fileUpload.FileType]
	file, err := os.Create(filePath)

	fmt.Println(filePath)

	if err != nil {
		http.Error(w, "Error creating filepath: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = file.Write(fileUpload.File)

	if err != nil {
		http.Error(w, "Error writing to file: "+err.Error(), http.StatusBadRequest)
		return
	}

	res := schema.Response{
		Message: "File saved successfully!",
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

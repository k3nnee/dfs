package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/k3nnee/dfs/backend/master_node/app/schema"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

const maxFileSize = 10 << 20
const blockSize = 5 << 20

var allowedFileTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func spliceData(file []byte, fileType string) (map[string]string, bool) {
	workerNodes := strings.Split(os.Getenv("WORKER_NODES"), ",")
	randArr := make([]int, len(workerNodes))

	for index, _ := range workerNodes {
		randArr[index] = index
	}

	rand.Shuffle(len(randArr), func(i, j int) {
		randArr[i], randArr[j] = randArr[j], randArr[i]
	})

	var splitData [][]byte

	if len(file) > blockSize {
		splitData = append(splitData, file[:len(file)/2])
		splitData = append(splitData, file[len(file)/2:])
	} else {
		splitData = append(splitData, file)
	}

	metadata := make(map[string]string)
	for index, b := range splitData {
		hash := sha1.Sum(b)
		metadata[hex.EncodeToString(hash[:])] = workerNodes[randArr[index]]

		payload := schema.FileUpload{
			File:     b,
			FileName: hex.EncodeToString(hash[:]),
			FileType: fileType,
		}

		payloadJSON, err := json.Marshal(payload)

		if err != nil {
			return nil, false
		}

		res, err := http.Post(workerNodes[randArr[index]], "application/json", bytes.NewBuffer(payloadJSON))

		if err != nil || res.StatusCode != http.StatusOK {
			if res != nil {
				body, _ := io.ReadAll(res.Body)
				res.Body.Close()
				fmt.Println(string(body))
			}
			fmt.Println(err.Error())
			return nil, false
		}

		res.Body.Close()

		metadata[hex.EncodeToString(hash[:])+"BACKUP"] = workerNodes[randArr[index]]

		payload = schema.FileUpload{
			File:     b,
			FileName: hex.EncodeToString(hash[:]) + "BACKUP",
			FileType: fileType,
		}

		payloadJSON, err = json.Marshal(payload)

		if err != nil {
			return nil, false
		}

		res, err = http.Post(workerNodes[randArr[(index+1)%len(workerNodes)]], "application/json", bytes.NewBuffer(payloadJSON))

		if err != nil || res.StatusCode != http.StatusOK {
			if res != nil {
				body, _ := io.ReadAll(res.Body)
				res.Body.Close()
				fmt.Println(string(body))
			}
			fmt.Println(err.Error())
			return nil, false
		}

		res.Body.Close()
	}

	return metadata, true
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
		http.Error(w, "File too large, limit of 10mb", http.StatusRequestEntityTooLarge)
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

	fileContent, err := io.ReadAll(file)

	if err != nil {
		http.Error(w, "Unable to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, ok := spliceData(append(buf, fileContent...), fileType)

	if !ok {
		http.Error(w, "Error partitioning data", http.StatusInternalServerError)
		return
	}

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

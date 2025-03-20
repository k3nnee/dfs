package main

import (
	"encoding/json"
	"fmt"
	"github.com/k3nnee/dfs/backend/worker_node/app/schema"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")

		w.Write([]byte("hello world"))
	})

	handler.HandleFunc("/upload-data", func(w http.ResponseWriter, r *http.Request) {
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

		fmt.Println(fileUpload)

	})

	s := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}

package main

import (
	api "github.com/k3nnee/dfs/backend/worker_node/app/routes"
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

	handler.HandleFunc("/upload-data", api.UploadData)

	s := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}

package main

import (
	api "github.com/k3nnee/dfs/backend/master_node/app/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	handler := http.NewServeMux()

	handler.HandleFunc("/upload-data", api.UploadData)

	s := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}

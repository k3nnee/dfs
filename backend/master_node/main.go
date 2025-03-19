package main

import (
	api "github.com/k3nnee/dfs/backend/master_node/app/routes"
	"log"
	"net/http"
)

func main() {
	handler := http.NewServeMux()

	handler.HandleFunc("/upload-data", api.UploadData)

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}

package main

import (
	database "github.com/k3nnee/dfs/backend/master_node/app/db"
	api "github.com/k3nnee/dfs/backend/master_node/app/routes"
	"log"
	"net/http"
)

func main() {
	_ = database.DBConfig()
	handler := http.NewServeMux()

	handler.HandleFunc("/upload-data", api.UploadData)

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}

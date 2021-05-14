package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed client/dist
var clientFS embed.FS

func main() {
	distFS, err := fs.Sub(clientFS, "client/dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(distFS)))

	log.Println("Starting HTTP server at http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

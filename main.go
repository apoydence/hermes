package main

import (
	"github.com/apoydence/hermes/handlers"
	"net/http"
)

func main() {
	keyStorer := handlers.NewKeyStorer()
	uploader := handlers.NewUploader(keyStorer)
	downloader := handlers.NewDownloader(keyStorer)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if req.Method == "GET" {
			downloader.ServeHTTP(w, req)
		} else if req.Method == "POST" {
			uploader.ServeHTTP(w, req)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

package handlers

import (
	"io"
	"net/http"
)

type KeyFetcher interface {
	Fetch(key string) io.Reader
}

type Downloader struct {
	keyFetcher KeyFetcher
}

func NewDownloader(keyFetcher KeyFetcher) *Downloader {
	return &Downloader{
		keyFetcher: keyFetcher,
	}
}

func (d *Downloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("key")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key header is required"))
		return
	}

	reader := d.keyFetcher.Fetch(cookie.Value)
	if reader == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			return
		}
		w.Write(buf[:n])
	}
}

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
	key := r.Header.Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key header is required"))
		return
	}

	reader := d.keyFetcher.Fetch(key)
	if reader == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			return
		}
		w.Write(buf[:n])
	}
}

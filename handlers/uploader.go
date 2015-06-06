package handlers

import (
	"io"
	"net/http"
)

type KeyStorage interface {
	Add(key string, data io.ReadCloser) error
}

type Uploader struct {
	keyStorage KeyStorage
}

func NewUploader(keyStorage KeyStorage) *Uploader {
	return &Uploader{
		keyStorage: keyStorage,
	}
}

func (u *Uploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key header is required"))
		return
	}

	err := u.keyStorage.Add(key, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
}

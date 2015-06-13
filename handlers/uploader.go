package handlers

import (
	"io"
	"net/http"
	"time"
)

type KeyStorage interface {
	Add(key string, data io.Reader) error
	Delete(key string)
}

type Uploader struct {
	keyStorage KeyStorage
	timeout    time.Duration
}

func NewUploader(keyStorage KeyStorage, timeout time.Duration) *Uploader {
	return &Uploader{
		keyStorage: keyStorage,
		timeout:    timeout,
	}
}

func (u *Uploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key header is required"))
		return
	}

	chReader := NewChannelReader(r.Body)
	err := u.keyStorage.Add(key, chReader)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
	defer u.keyStorage.Delete(key)
	defer r.Body.Close()
	chReader.Run(u.timeout)
}

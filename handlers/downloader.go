package handlers

import (
	"io"
	"net/http"
)

type KeyFetcher interface {
	Fetch(key string) *ReadLocker
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
	notifier, ok := w.(http.CloseNotifier)
	if !ok {
		panic("ResponseWriter is not a CloseNotifier")
	}
	notify := notifier.CloseNotify()

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

	if ok := reader.Lock.TryLock(); !ok {
		w.WriteHeader(http.StatusConflict)
		return
	}
	defer reader.Lock.Unlock()

	w.Header().Set("Content-Type", "application/octet-stream")
	dataCh := run(reader.Reader, notify)
	for {
		select {
		case <-notify:
			return
		case buf, ok := <-dataCh:
			if !ok {
				return
			}
			w.Write(buf)
		}
	}
}

func run(reader io.Reader, notify <-chan bool) <-chan []byte {
	c := make(chan []byte)
	go func() {
		defer close(c)
		buf := make([]byte, 1024)
		for !isDone(notify) {
			n, err := reader.Read(buf)
			if err != nil {
				return
			}
			newBuf := make([]byte, n)
			for i := 0; i < n; i++ {
				newBuf[i] = buf[i]
			}
			c <- newBuf
		}
	}()
	return c
}

func isDone(c <-chan bool) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}

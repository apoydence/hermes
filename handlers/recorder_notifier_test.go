package handlers_test

import "net/http/httptest"

type recorderNotifier struct {
	*httptest.ResponseRecorder
	notifyCh chan bool
}

func newRecorderNotifier() *recorderNotifier {
	r := &recorderNotifier{
		notifyCh: make(chan bool, 100),
	}
	r.ResponseRecorder = httptest.NewRecorder()
	return r
}

func (r *recorderNotifier) CloseNotify() <-chan bool {
	return r.notifyCh
}

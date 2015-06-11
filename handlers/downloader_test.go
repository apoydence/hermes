package handlers_test

import (
	"bytes"
	"github.com/apoydence/hermes/handlers"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Downloader", func() {
	var (
		downloader     *handlers.Downloader
		recorder       *httptest.ResponseRecorder
		mockKeyStorage mockKeyStorage
		req            *http.Request
		keyName        string
		cookie         *http.Cookie
	)

	BeforeEach(func() {
		keyName = "some-key"
		mockKeyStorage = NewMockKeyStorage()
		downloader = handlers.NewDownloader(mockKeyStorage)
		recorder = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "http://some.com", nil)
		cookie = &http.Cookie{
			Name:  "key",
			Value: keyName,
		}
	})

	Context("Happy", func() {
		It("Passes data through", func() {
			req.AddCookie(cookie)
			expectedData := []byte("here is some fake data")
			buf := bytes.NewBuffer(expectedData)
			mockKeyStorage.Add(keyName, buf)
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			contentType := recorder.HeaderMap.Get("Content-Type")
			Expect(contentType).To(Equal("application/octet-stream"))
			Expect(recorder.Body.Bytes()).To(Equal(expectedData))
		})
	})

	Context("Unhappy", func() {
		It("Returns StatusBadRequest if a key is not provided", func() {
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			Expect(recorder.Body.Bytes()).To(Equal([]byte("Key header is required")))
		})

		It("Returns StatusNotFound if the provided key is not found", func() {
			req.AddCookie(cookie)
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})

		It("Returns a StatusConflict if the key is already being read from", func() {
			req.AddCookie(cookie)
			expectedData := []byte("here is some fake data")
			buf := bytes.NewBuffer(expectedData)
			mockKeyStorage.Add(keyName, buf)
			ok := mockKeyStorage.Fetch(keyName).Lock.TryLock()

			Expect(ok).To(BeTrue())
			downloader.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusConflict))
		})
	})

})

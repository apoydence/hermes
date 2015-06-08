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
	)

	BeforeEach(func() {
		keyName = "some-key"
		mockKeyStorage = NewMockKeyStorage()
		downloader = handlers.NewDownloader(mockKeyStorage)
		recorder = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "http://some.com", nil)
		req.Header.Add("key", keyName)
	})

	Context("Happy", func() {
		It("Passes data through", func() {
			expectedData := []byte("here is some fake data")
			buf := bytes.NewBuffer(expectedData)
			mockKeyStorage[keyName] = buf
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Body.Bytes()).To(Equal(expectedData))
		})
	})

	Context("Unhappy", func() {
		It("Returns StatusBadRequest if a key is not provided", func() {
			req.Header.Del("key")
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			Expect(recorder.Body.Bytes()).To(Equal([]byte("Key header is required")))
		})

		It("Returns StatusNotFound if the provided key is not found", func() {
			req.Header.Add("key", keyName)
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})
	})

})

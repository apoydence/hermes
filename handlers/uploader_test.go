package handlers_test

import (
	"bytes"
	"github.com/apoydence/hermes/handlers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uploader", func() {

	var (
		uploader       *handlers.Uploader
		recorder       *httptest.ResponseRecorder
		mockKeyStorage mockKeyStorage
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		mockKeyStorage = NewMockKeyStorage()
		uploader = handlers.NewUploader(mockKeyStorage)
	})

	Context("Happy", func() {
		It("Data passes through", func() {
			expectedData := []byte("Here is some data that is expected to flow through")
			buf := bytes.NewBuffer(expectedData)
			req, err := http.NewRequest("POST", "http://some.com", buf)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Add("key", "some-key")
			uploader.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			getReader, ok := mockKeyStorage["some-key"]
			Expect(ok).To(BeTrue())
			Expect(getReader).ToNot(BeNil())
			data, err := ioutil.ReadAll(getReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(Equal(expectedData))
		})
	})

	Context("Unhappy", func() {
		It("Returns a StatusBadRequest if the key is not provided", func() {
			req, err := http.NewRequest("POST", "http://some.com", nil)
			Expect(err).ToNot(HaveOccurred())
			uploader.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			Expect(recorder.Body.Bytes()).To(Equal([]byte("Key header is required")))
		})

		It("Returns a StatusConflict if the key is taken", func() {
			mockKeyStorage["some-key"] = ioutil.NopCloser(&bytes.Buffer{})
			req, err := http.NewRequest("POST", "http://some.com", nil)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Add("key", "some-key")
			uploader.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusConflict))
			Expect(recorder.Body.Bytes()).To(Equal([]byte("Key already present")))
		})
	})

})

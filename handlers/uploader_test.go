package handlers_test

import (
	"bytes"
	"github.com/apoydence/hermes/handlers"
	"io"
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
		mockKeyStorage *handlers.KeyStorer
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		mockKeyStorage = handlers.NewKeyStorer()
		uploader = handlers.NewUploader(mockKeyStorage)
	})

	Context("Happy", func() {
		It("Passes data through", func(done Done) {
			defer close(done)
			expectedData := []byte("Here is some data that is expected to flow through")
			buf := bytes.NewBuffer(expectedData)
			reader := NewMockReadCloser(buf, 1)
			req, err := http.NewRequest("POST", "http://some.com", reader)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Add("key", "some-key")
			serveDone := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer close(serveDone)
				uploader.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusOK))
			}()
			var getReader io.Reader
			getReaderFunc := func() io.Reader {
				getReader = mockKeyStorage.Fetch("some-key")
				return getReader
			}
			Eventually(getReaderFunc).ShouldNot(BeNil())
			Consistently(serveDone).ShouldNot(BeClosed())
			data, err := ioutil.ReadAll(getReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(Equal(expectedData))
			Expect(reader.isClosed).To(BeTrue())
			Eventually(serveDone).Should(BeClosed())
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
			mockKeyStorage.Add("some-key", &bytes.Buffer{})
			req, err := http.NewRequest("POST", "http://some.com", nil)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Add("key", "some-key")
			uploader.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusConflict))
			Expect(recorder.Body.Bytes()).To(Equal([]byte("The key some-key is already in use")))
		})
	})

})

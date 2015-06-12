package handlers_test

import (
	"bytes"
	"github.com/apoydence/hermes/handlers"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Downloader", func() {
	var (
		downloader     *handlers.Downloader
		recorder       *recorderNotifier
		mockKeyStorage mockKeyStorage
		req            *http.Request
		keyName        string
		cookie         *http.Cookie
	)

	BeforeEach(func() {
		keyName = "some-key"
		mockKeyStorage = NewMockKeyStorage()
		downloader = handlers.NewDownloader(mockKeyStorage)
		recorder = newRecorderNotifier()
		req, _ = http.NewRequest("GET", "http://some.com", nil)
		cookie = &http.Cookie{
			Name:  "key",
			Value: keyName,
		}
	})

	Context("Happy", func() {
		It("Passes data through", func() {
			req.AddCookie(cookie)
			var expectedData []byte
			for i := 0; i < 2048; i++ {
				expectedData = append(expectedData, byte(i))
			}

			buf := bytes.NewBuffer(expectedData)
			mockKeyStorage.Add(keyName, buf)
			downloader.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			contentType := recorder.HeaderMap.Get("Content-Type")
			Expect(contentType).To(Equal("application/octet-stream"))
			Expect(recorder.Body.Bytes()).To(Equal(expectedData))
		})

		It("Stops downloading if the notify channel is written to", func() {
			req.AddCookie(cookie)
			expectedData := []byte("here is some fake data")
			buf := bytes.NewBuffer(expectedData)
			chReader := handlers.NewChannelReader(ioutil.NopCloser(buf))
			mockKeyStorage.Add(keyName, chReader)
			doneCh := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer close(doneCh)
				downloader.ServeHTTP(recorder, req)
			}()

			recorder.notifyCh <- true
			Eventually(doneCh).Should(BeClosed())
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

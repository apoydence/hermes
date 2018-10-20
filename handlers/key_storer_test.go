package handlers_test

import (
	"bytes"
	"github.com/poy/hermes/handlers"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyStorage", func() {
	var (
		keyStorer *handlers.KeyStorer
	)

	BeforeEach(func() {
		keyStorer = handlers.NewKeyStorer()
	})

	It("returns the reader for the given key", func() {
		buf := &bytes.Buffer{}
		go func() {
			defer GinkgoRecover()
			err := keyStorer.Add("some-key", buf)
			Expect(err).ToNot(HaveOccurred())
		}()

		f := func() io.Reader {
			r := keyStorer.Fetch("some-key")
			if r != nil {
				return r.Reader
			}
			return nil
		}
		Eventually(f).Should(Equal(buf))
	})

	It("returns nil for an unknown key", func() {
		reader := keyStorer.Fetch("some-key")
		Expect(reader).To(BeNil())
	})

	It("returns an error if the same key is added twice", func() {
		buf := &bytes.Buffer{}
		err := keyStorer.Add("some-key", buf)
		Expect(err).ToNot(HaveOccurred())
		err = keyStorer.Add("some-key", buf)
		Expect(err).To(HaveOccurred())
	})

})

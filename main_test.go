package main_test

import (
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	It("compiles", func() {
		_, err := gexec.Build("github.com/apoydence/hermes")
		Expect(err).ToNot(HaveOccurred())
	})
})

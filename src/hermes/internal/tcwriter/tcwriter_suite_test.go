package tcwriter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTcwriter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tcwriter Suite")
}

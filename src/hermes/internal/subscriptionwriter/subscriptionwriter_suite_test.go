package subscriptionwriter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSubscriptionwriter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subscriptionwriter Suite")
}

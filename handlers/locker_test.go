package handlers_test

import (
	"github.com/apoydence/hermes/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Locker", func() {
	var (
		locker *handlers.Locker
	)

	BeforeEach(func() {
		locker = handlers.NewLocker()
	})

	It("can only acquire a lock one at a time", func() {
		ok := locker.TryLock()
		Expect(ok).To(BeTrue())
		ok = locker.TryLock()
		Expect(ok).To(BeFalse())
		locker.Unlock()
		ok = locker.TryLock()
		Expect(ok).To(BeTrue())
	})

})

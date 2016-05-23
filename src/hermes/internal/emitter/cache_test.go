package emitter_test

import (
	"hermes/internal/datastructures"
	"hermes/internal/emitter"

	. "github.com/apoydence/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var (
		mockRegistry *mockRegistry
		expectedList *datastructures.LinkedList

		senderStores chan emitter.SenderStore

		cache      *emitter.Cache
		expectedId string
	)

	var origSubscriptionListReader = emitter.NewSubscriptionListReader
	var hijackedSubscriptionListReader = func(senderStore emitter.SenderStore) *emitter.SubscriptionListReader {
		senderStores <- senderStore
		return origSubscriptionListReader(senderStore)
	}

	BeforeEach(func() {
		mockRegistry = newMockRegistry()
		expectedList = new(datastructures.LinkedList)
		senderStores = make(chan emitter.SenderStore, 100)

		emitter.NewSubscriptionListReader = hijackedSubscriptionListReader

		cache = emitter.NewCache(mockRegistry)

		expectedId = "some-id"
	})

	AfterEach(func() {
		emitter.NewSubscriptionListReader = origSubscriptionListReader
	})

	Describe("Fetch()", func() {
		JustBeforeEach(func() {
			close(mockRegistry.GetListOutput.Ret0)
		})

		Context("registry returns list", func() {
			BeforeEach(func() {
				mockRegistry.GetListOutput.Ret0 <- expectedList
			})

			Context("empty cache", func() {
				It("requests a the list from the registry", func() {
					cache.Fetch(expectedId)

					Expect(mockRegistry.GetListInput).To(BeCalled(With(expectedId)))
				})

				It("returns the expected list", func() {
					cache.Fetch(expectedId)

					Expect(senderStores).To(Receive(Equal(expectedList)))
				})

				Context("cache has entry", func() {
					BeforeEach(func() {
						cache.Fetch(expectedId)
					})

					It("only uses factory once for a given ID", func() {
						cache.Fetch(expectedId)

						Expect(mockRegistry.GetListCalled).To(HaveLen(1))
					})
				})
			})
		})
	})

})

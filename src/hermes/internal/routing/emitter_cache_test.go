package routing_test

import (
	"hermes/internal/routing"

	. "github.com/apoydence/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmitterCache", func() {
	var (
		mockEmitterFetcher *mockEmitterFetcher
		mockEmitter        *mockEmitter

		cache      *routing.EmitterCache
		expectedId string
	)

	BeforeEach(func() {
		mockEmitterFetcher = newMockEmitterFetcher()
		mockEmitter = newMockEmitter()

		cache = routing.NewEmitterCache(mockEmitterFetcher)

		expectedId = "some-id"
	})

	Describe("Fetch()", func() {
		JustBeforeEach(func() {
			close(mockEmitterFetcher.FetchOutput.Ret0)
		})

		Context("factory returns emitter", func() {
			BeforeEach(func() {
				mockEmitterFetcher.FetchOutput.Ret0 <- mockEmitter
			})

			Context("empty cache", func() {
				It("requests a new emitter from the factory", func() {
					cache.Fetch(expectedId)

					Expect(mockEmitterFetcher.FetchInput).To(BeCalled(With(expectedId)))
				})

				It("returns the expected emitter", func() {
					emitter := cache.Fetch(expectedId)

					Expect(emitter).To(Equal(mockEmitter))
				})

				Context("cache has entry", func() {
					BeforeEach(func() {
						cache.Fetch(expectedId)
					})

					It("only uses factory once for a given ID", func() {
						cache.Fetch(expectedId)

						Expect(mockEmitterFetcher.FetchCalled).To(HaveLen(1))
					})
				})
			})
		})
	})

})

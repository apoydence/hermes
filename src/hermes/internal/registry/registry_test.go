//go:generate hel
package registry_test

import (
	"fmt"
	"hermes/internal/datastructures"
	"hermes/internal/emitter"
	"hermes/internal/registry"
	"hermes/internal/subscriptionwriter"
	"unsafe"

	. "github.com/apoydence/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {
	var (
		linkedLists            chan *datastructures.LinkedList
		subscriptionWriters    chan *subscriptionwriter.SubscriptionWriter
		subscriptionWriterUrls chan string
		subscriptionWriterErrs chan error
		muxIds                 chan uint64
		msgEmitters            chan subscriptionwriter.MessageEmitter
		mockKvStore            *mockKvStore

		reg *registry.Registry
	)

	var origLinkedList = datastructures.NewLinkedList
	var hijackLinkedList = func() *datastructures.LinkedList {
		result := origLinkedList()
		linkedLists <- result
		return result
	}

	var origNewEmitter = subscriptionwriter.NewEmitter
	var hijackedTcEmitter = func(muxId uint64, msgEmitter subscriptionwriter.MessageEmitter) *subscriptionwriter.Emitter {
		muxIds <- muxId
		msgEmitters <- msgEmitter

		return origNewEmitter(muxId, msgEmitter)
	}

	var origNewSubscriptionWriter = subscriptionwriter.New
	var hijackedNewSubscriptionWriter = func(url string) (*subscriptionwriter.SubscriptionWriter, error) {
		subscriptionWriterUrls <- url
		writer := new(subscriptionwriter.SubscriptionWriter)
		subscriptionWriters <- writer
		return writer, <-subscriptionWriterErrs
	}

	BeforeEach(func() {
		datastructures.NewLinkedList = hijackLinkedList
		subscriptionwriter.NewEmitter = hijackedTcEmitter
		subscriptionwriter.New = hijackedNewSubscriptionWriter

		linkedLists = make(chan *datastructures.LinkedList, 100)
		muxIds = make(chan uint64, 100)
		msgEmitters = make(chan subscriptionwriter.MessageEmitter, 100)
		subscriptionWriters = make(chan *subscriptionwriter.SubscriptionWriter, 100)
		subscriptionWriterUrls = make(chan string, 100)
		subscriptionWriterErrs = make(chan error, 100)

		mockKvStore = newMockKvStore()

		reg = registry.New(mockKvStore)
	})

	AfterEach(func() {
		datastructures.NewLinkedList = origLinkedList
		subscriptionwriter.NewEmitter = origNewEmitter
		subscriptionwriter.New = origNewSubscriptionWriter
	})

	Describe("GetList()", func() {
		var (
			expectedId string
		)

		BeforeEach(func() {
			expectedId = "some-id"
		})

		Context("list is not yet created", func() {
			It("creates a new list", func() {
				reg.GetList(expectedId)

				Expect(linkedLists).To(HaveLen(1))
			})
		})

		Context("list has already been requested", func() {
			var (
				origList *datastructures.LinkedList
			)

			BeforeEach(func() {
				origList = reg.GetList(expectedId)
			})

			It("does not create a new list", func() {
				reg.GetList(expectedId)

				Expect(linkedLists).To(HaveLen(1))
			})

			It("returns the same LinkedList", func() {
				list := reg.GetList(expectedId)

				Expect(list).To(Equal(origList))
			})

			It("creates a new list for a different ID", func() {
				reg.GetList(expectedId + "-other")

				Expect(linkedLists).To(HaveLen(2))
			})
		})

		Describe("KvStore manipulation", func() {
			var (
				callback      func(id, URL, key string, muxId uint64, add bool)
				expectedId    string
				expectedUrl   string
				expectedKey   string
				expectedMuxId uint64
			)

			BeforeEach(func() {
				expectedId = "some-id"
				expectedUrl = "some-url"
				expectedKey = "some-key"
				expectedMuxId = 99

				Expect(mockKvStore.ListenForInput.Callback).To(Receive(&callback))
			})

			Context("subscriptionWriter construction does not return an error", func() {
				BeforeEach(func() {
					close(subscriptionWriterErrs)
				})

				Context("KvStore adds subscription", func() {
					BeforeEach(func() {
						callback(expectedId, expectedUrl, expectedKey, expectedMuxId, true)
					})

					Describe("LinkedList interaction", func() {
						var fetchEmitters = func() []*emitter.Emitter {
							var emitters []*emitter.Emitter
							var list *datastructures.LinkedList
							Expect(linkedLists).To(Receive(&list))

							list.Traverse(func(value unsafe.Pointer) {
								emitters = append(emitters, (*emitter.Emitter)(value))
							})

							return emitters
						}

						It("creates a new LinkedList", func() {
							Expect(linkedLists).To(HaveLen(1))
						})

						It("adds the Emitter to the LinkedList", func() {
							Expect(fetchEmitters()).To(HaveLen(1))
						})
					})

					Describe("SubscriptionWriter construction", func() {
						It("creates a new SubscriptionWriter", func() {
							Expect(subscriptionWriterUrls).To(HaveLen(1))
						})

						It("creates a new SubscriptionWriter with the expected URL", func() {
							Expect(subscriptionWriterUrls).To(Receive(Equal(expectedUrl)))
						})
					})

					Describe("Emitter construction", func() {
						It("creates a new Emitter", func() {
							Expect(msgEmitters).To(HaveLen(1))
						})

						It("creates a new Emitter with the expected MuxID", func() {
							Expect(muxIds).To(Receive(Equal(expectedMuxId)))
						})

						It("creates a new Emitter with the expected MuxID", func() {
							var subscriptionWriter *subscriptionwriter.SubscriptionWriter
							Expect(subscriptionWriters).To(Receive(&subscriptionWriter))

							Expect(msgEmitters).To(Receive(Equal(subscriptionWriter)))
						})
					})
				})
			})

			Context("subscriptionWriter construction returns an error", func() {
				var (
					expectedSubscriptionWriterErr error
				)

				BeforeEach(func() {
					expectedSubscriptionWriterErr = fmt.Errorf("some-error")
					subscriptionWriterErrs <- expectedSubscriptionWriterErr
					callback(expectedId, expectedUrl, expectedKey, expectedMuxId, true)
				})

				It("removes the key", func() {
					Expect(mockKvStore.RemoveInput).To(BeCalled(With(expectedKey)))
				})
			})
		})
	})
})

package emitter_test

import (
	"hermes/internal/datastructures"
	"hermes/internal/emitter"
	"hermes/internal/tcwriter"
	"unsafe"

	. "github.com/apoydence/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	var (
		generator *emitter.Generator

		linkedLists   chan *datastructures.LinkedList
		senderStores  chan emitter.SenderStore
		eventEmitters chan *emitter.SubscriptionListReader
		muxIds        chan uint64
		msgEmitters   chan tcwriter.MessageEmitter
		tcwriterUrls  chan string
		tcWriterErrs  chan error

		mockKvStore *mockKvStore

		expectedId string
	)

	var origLinkedList = datastructures.NewLinkedList
	var hijackLinkedList = func() *datastructures.LinkedList {
		result := origLinkedList()
		linkedLists <- result
		return result
	}

	var origSubscriptionListReader = emitter.NewSubscriptionListReader
	var hijackedSubscriptionListReader = func(store emitter.SenderStore) *emitter.SubscriptionListReader {
		senderStores <- store
		result := origSubscriptionListReader(store)
		eventEmitters <- result
		return result
	}

	var origNewTcEmitter = tcwriter.NewEmitter
	var hijackedTcEmitter = func(muxId uint64, msgEmitter tcwriter.MessageEmitter) *tcwriter.Emitter {
		muxIds <- muxId
		msgEmitters <- msgEmitter

		return origNewTcEmitter(muxId, msgEmitter)
	}

	var origNewTcWriter = tcwriter.New
	var hijackedNewTcWriter = func(url string) (*tcwriter.TcWriter, error) {
		tcwriterUrls <- url
		return new(tcwriter.TcWriter), <-tcWriterErrs
	}

	BeforeEach(func() {
		datastructures.NewLinkedList = hijackLinkedList
		emitter.NewSubscriptionListReader = hijackedSubscriptionListReader
		tcwriter.NewEmitter = hijackedTcEmitter
		tcwriter.New = hijackedNewTcWriter

		mockKvStore = newMockKvStore()

		generator = emitter.NewGenerator(mockKvStore)

		linkedLists = make(chan *datastructures.LinkedList, 100)
		senderStores = make(chan emitter.SenderStore, 100)
		eventEmitters = make(chan *emitter.SubscriptionListReader, 100)
		muxIds = make(chan uint64, 100)
		msgEmitters = make(chan tcwriter.MessageEmitter, 100)
		tcwriterUrls = make(chan string, 100)
		tcWriterErrs = make(chan error, 100)

		expectedId = "some-id"
	})

	AfterEach(func() {
		datastructures.NewLinkedList = origLinkedList
		emitter.NewSubscriptionListReader = origSubscriptionListReader
		tcwriter.NewEmitter = origNewTcEmitter
		tcwriter.New = origNewTcWriter
	})

	Describe("Fetch()", func() {
		Context("tcwriter construction does not return an error", func() {
			BeforeEach(func() {
				close(tcWriterErrs)
			})
			It("creates a new linked list", func() {
				generator.Fetch(expectedId)

				Expect(linkedLists).To(HaveLen(1))
			})

			It("creates a new SubscriptionListReader with the linked list", func() {
				generator.Fetch(expectedId)

				var ll *datastructures.LinkedList
				Expect(linkedLists).To(Receive(&ll))
				Expect(senderStores).To(BeCalled(With(ll)))
			})

			It("returns the SubscriptionListReader", func() {
				result := generator.Fetch(expectedId)

				var ee *emitter.SubscriptionListReader
				Expect(eventEmitters).To(Receive(&ee))
				Expect(result).To(Equal(ee))
			})
		})
	})

	Describe("KvStore interaction", func() {
		Context("tcwriter construction does not return an error", func() {
			var (
				callback   func(URL string, muxId uint64)
				linkedList *datastructures.LinkedList

				expectedURL   string
				expectedMuxId uint64
			)

			var listLen = func(ll *datastructures.LinkedList) int {
				var count int
				ll.Traverse(func(unsafe.Pointer) {
					count++
				})

				return count
			}

			var fetchCallback = func() {
				Expect(mockKvStore.ListenForInput.Callback).To(Receive(&callback))
			}

			var fetchLinkedList = func() {
				Expect(linkedLists).To(Receive(&linkedList))
			}

			BeforeEach(func() {
				expectedURL = "http://some.url"
				expectedMuxId = 99

				close(tcWriterErrs)
				generator.Fetch(expectedId)
				fetchLinkedList()
				fetchCallback()
			})

			Context("nothing added", func() {

				It("doesn't add anything to the linked list", func() {
					Expect(listLen(linkedList)).To(Equal(0))
				})

				It("listens to the expected id", func() {
					Expect(mockKvStore.ListenForInput.Id).To(BeCalled(With(expectedId)))
				})

				Context("adds listener", func() {
					BeforeEach(func() {
						callback(expectedURL, expectedMuxId)
					})

					It("adds a single entry", func() {
						Expect(listLen(linkedList)).To(Equal(1))
					})

					Describe("NewEmitter construction", func() {
						It("builds the emitter with the expected mux id", func() {
							Expect(muxIds).To(Receive(Equal(expectedMuxId)))
						})

						It("uses a non-nil emitter", func() {
							var emitter tcwriter.MessageEmitter
							Expect(msgEmitters).To(Receive(&emitter))
							Expect(emitter).ToNot(BeNil())
						})
					})

					Describe("TcWriter construction", func() {
						It("builds the TcWriter with the expected URL", func() {
							Expect(tcwriterUrls).To(Receive(Equal(expectedURL)))
						})
					})
				})
			})
		})
	})

})

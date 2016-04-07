package routing_test

import (
	"hermes/internal/routing"
	"hermes/internal/tcwriter"

	. "github.com/apoydence/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmitterGenerator", func() {
	var (
		generator *routing.EmitterGenerator

		linkedLists   chan *routing.LinkedList
		senderStores  chan routing.SenderStore
		eventEmitters chan *routing.EventEmitter
		muxIds        chan uint64
		msgEmitters   chan tcwriter.MessageEmitter
		tcwriterUrls  chan string
		tcWriterErrs  chan error

		mockKvStore *mockKvStore

		expectedId string
	)

	var origLinkedList = routing.NewLinkedList
	var hijackLinkedList = func() *routing.LinkedList {
		result := origLinkedList()
		linkedLists <- result
		return result
	}

	var origEventEmitter = routing.NewEventEmitter
	var hijackedEventEmitter = func(store routing.SenderStore) *routing.EventEmitter {
		senderStores <- store
		result := origEventEmitter(store)
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
		routing.NewLinkedList = hijackLinkedList
		routing.NewEventEmitter = hijackedEventEmitter
		tcwriter.NewEmitter = hijackedTcEmitter
		tcwriter.New = hijackedNewTcWriter

		mockKvStore = newMockKvStore()

		generator = routing.NewEmitterGenerator(mockKvStore)

		linkedLists = make(chan *routing.LinkedList, 100)
		senderStores = make(chan routing.SenderStore, 100)
		eventEmitters = make(chan *routing.EventEmitter, 100)
		muxIds = make(chan uint64, 100)
		msgEmitters = make(chan tcwriter.MessageEmitter, 100)
		tcwriterUrls = make(chan string, 100)
		tcWriterErrs = make(chan error, 100)

		expectedId = "some-id"
	})

	AfterEach(func() {
		routing.NewLinkedList = origLinkedList
		routing.NewEventEmitter = origEventEmitter
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

			It("creates a new EventEmitter with the linked list", func() {
				generator.Fetch(expectedId)

				var ll *routing.LinkedList
				Expect(linkedLists).To(Receive(&ll))
				Expect(senderStores).To(BeCalled(With(ll)))
			})

			It("returns the EventEmitter", func() {
				result := generator.Fetch(expectedId)

				var ee *routing.EventEmitter
				Expect(eventEmitters).To(Receive(&ee))
				Expect(result).To(Equal(ee))
			})
		})
	})

	Describe("KvStore interaction", func() {
		Context("tcwriter construction does not return an error", func() {
			var (
				callback   func(URL string, muxId uint64)
				linkedList *routing.LinkedList

				expectedURL   string
				expectedMuxId uint64
			)

			var listLen = func(ll *routing.LinkedList) int {
				var count int
				ll.Traverse(func(routing.Emitter) {
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

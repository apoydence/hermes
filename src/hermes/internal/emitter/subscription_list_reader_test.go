//go:generate hel
package emitter_test

import (
	"hermes/common/pb/messages"
	"hermes/internal/emitter"
	"unsafe"

	. "github.com/apoydence/eachers"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StoreReader", func() {
	var (
		mockSenderStore    *mockSenderStore
		mockEmitter        *mockEmitter
		subscriptionReader *emitter.SubscriptionListReader
		expextedData       *messages.DataPoint
	)

	BeforeEach(func() {
		mockSenderStore = newMockSenderStore()
		mockEmitter = newMockEmitter()
		subscriptionReader = emitter.NewSubscriptionListReader(mockSenderStore)
		expextedData = &messages.DataPoint{
			Id: proto.String("some-id"),
		}
	})

	Context("store reader does not return an error", func() {
		BeforeEach(func() {
			close(mockEmitter.EmitOutput.Ret0)
		})

		It("uses listener store", func() {
			subscriptionReader.Emit(expextedData)

			Expect(mockSenderStore.TraverseCalled).To(HaveLen(1))
		})

		It("sends each sender the expected data", func() {
			subscriptionReader.Emit(expextedData)

			var callback func(unsafe.Pointer)
			Expect(mockSenderStore.TraverseInput.Callback).To(Receive(&callback))
			var e emitter.Emitter = mockEmitter
			callback(unsafe.Pointer(&e))
			callback(unsafe.Pointer(&e))

			Expect(mockEmitter.EmitInput.Data).To(EqualEach(expextedData, expextedData))
		})
	})
})

//go:generate hel
package routing_test

import (
	"hermes/common/pb/messages"
	"hermes/doppler/internal/routing"

	. "github.com/apoydence/eachers"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventEmitter", func() {
	var (
		mockSenderStore *mockSenderStore
		mockEmitter     *mockEmitter
		emitter         *routing.EventEmitter
		expextedData    *messages.DataPoint
	)

	BeforeEach(func() {
		mockSenderStore = newMockSenderStore()
		mockEmitter = newMockEmitter()
		emitter = routing.NewEventEmitter(mockSenderStore)
		expextedData = &messages.DataPoint{
			Id: proto.String("some-id"),
		}
	})

	Context("emitter does not return an error", func() {
		BeforeEach(func() {
			close(mockEmitter.EmitOutput.Ret0)
		})

		It("uses listener store", func() {
			emitter.Emit(expextedData)

			Expect(mockSenderStore.TraverseCalled).To(HaveLen(1))
		})

		It("sends each sender the expected data", func() {
			emitter.Emit(expextedData)

			var callback func(routing.Emitter)
			Expect(mockSenderStore.TraverseInput.Callback).To(Receive(&callback))
			callback(mockEmitter)
			callback(mockEmitter)

			Expect(mockEmitter.EmitInput.Data).To(EqualEach(expextedData, expextedData))
		})
	})
})

//go:generate hel
package tcwriter_test

import (
	"fmt"
	"hermes/common/pb/messages"
	"hermes/doppler/internal/tcwriter"

	. "github.com/apoydence/eachers"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Emitter", func() {
	var (
		mockMessageEmitter *mockMessageEmitter
		emitter            *tcwriter.Emitter

		expectedMessage *messages.Doppler
	)

	BeforeEach(func() {
		mockMessageEmitter = newMockMessageEmitter()

		expectedMessage = &messages.Doppler{
			MuxId:       proto.Uint64(99),
			MessageType: messages.Doppler_DataPoint.Enum(),
			DataPoint: &messages.DataPoint{
				Id: proto.String("some-id"),
			},
		}

		emitter = tcwriter.NewEmitter(expectedMessage.GetMuxId(), mockMessageEmitter)
	})

	Describe("Emit()", func() {
		Context("MessageEmitter does not return an error", func() {
			BeforeEach(func() {
				close(mockMessageEmitter.EmitMessageOutput.Ret0)
			})

			It("does not return an error", func() {
				Expect(emitter.Emit(expectedMessage.DataPoint)).To(Succeed())
			})

			It("encodes the expected MuxId in the message", func() {
				emitter.Emit(expectedMessage.DataPoint)

				Expect(mockMessageEmitter.EmitMessageInput).To(BeCalled(With(expectedMessage)))
			})
		})

		Context("MessageEmiter returns an error", func() {
			var (
				expectedErr error
			)

			BeforeEach(func() {
				expectedErr = fmt.Errorf("some-error")
				mockMessageEmitter.EmitMessageOutput.Ret0 <- expectedErr
			})

			It("returns an error", func() {
				Expect(emitter.Emit(expectedMessage.DataPoint)).To(MatchError(expectedErr))
			})
		})
	})
})

package routing_test

import (
	"hermes/common/pb/messages"
	"hermes/internal/routing"

	. "github.com/apoydence/eachers"
	. "github.com/apoydence/eachers/testhelpers"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var (
		mockDataSource     *mockDataSource
		mockWaitReporter   *mockWaitReporter
		mockEmitterFetcher *mockEmitterFetcher
		mockEmitter        *mockEmitter

		expectedDataPoint *messages.DataPoint
		expectedID        string

		router *routing.Router
	)

	BeforeEach(func() {
		mockDataSource = newMockDataSource()
		mockWaitReporter = newMockWaitReporter()
		mockEmitterFetcher = newMockEmitterFetcher()
		mockEmitter = newMockEmitter()

		expectedID = "some-id"
		expectedDataPoint = &messages.DataPoint{
			Id: proto.String(expectedID),
		}

		router = routing.New(mockDataSource, mockWaitReporter, mockEmitterFetcher)
	})

	Describe("Emit()", func() {
		JustBeforeEach(func() {
			AlwaysReturn(mockDataSource.NextOutput, new(messages.DataPoint), false)
			close(mockEmitterFetcher.FetchOutput.Ret0)
		})

		Context("no data ready", func() {
			It("reports that it is waiting", func() {
				Eventually(mockWaitReporter.WaitingCalled).ShouldNot(BeEmpty())
			})
		})
	})

	Context("data ready", func() {

		BeforeEach(func() {
			mockDataSource.NextOutput.Ret0 <- expectedDataPoint
			mockDataSource.NextOutput.Ret1 <- true
			mockEmitterFetcher.FetchOutput.Ret0 <- mockEmitter
		})

		Context("emitter does not return an error", func() {
			BeforeEach(func() {
				close(mockEmitter.EmitOutput.Ret0)
			})

			It("fetches an emitter for the ID and Waiting is not invoked", func() {
				Eventually(mockEmitterFetcher.FetchInput).Should(BeCalled(With(expectedID)))
				Expect(mockWaitReporter.WaitingCalled).To(BeEmpty())
			})

			It("emits data via the given emitter", func() {
				Eventually(mockEmitter.EmitInput).Should(BeCalled(With(expectedDataPoint)))
			})

			Context("multiple data points", func() {
				var (
					count int
				)

				BeforeEach(func() {
					count = 3
					for i := 0; i < count; i++ {
						mockDataSource.NextOutput.Ret0 <- expectedDataPoint
						mockDataSource.NextOutput.Ret1 <- true
						mockEmitterFetcher.FetchOutput.Ret0 <- mockEmitter
					}
				})

				It("fetches an emitter for each data point", func() {
					Eventually(mockEmitterFetcher.FetchCalled).Should(HaveLen(count + 1))
				})
			})
		})
	})

})

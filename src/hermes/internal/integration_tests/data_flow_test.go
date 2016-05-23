package integration_test

import (
	"fmt"
	"hermes/common/pb/messages"
	"hermes/internal/emitter"
	"hermes/internal/registry"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DataFlow", func() {
	var (
		mockServer *httptest.Server
		handlerWg  sync.WaitGroup

		mockDataSource   *mockDataSource
		mockWaitReporter *mockWaitReporter
		mockKvStore      *mockKvStore

		dataSourceReader *emitter.DataSourceReader
		cache            *emitter.Cache
		reg              *registry.Registry

		dopplerMessages chan *messages.Doppler
		conns           chan *websocket.Conn
	)

	var decodeMessage = func(data []byte) *messages.Doppler {
		msg := new(messages.Doppler)
		Expect(msg.Unmarshal(data)).To(Succeed())
		return msg
	}

	var convertHttpToWs = func(URL string) string {
		return fmt.Sprintf("ws%s", URL[4:])
	}

	var handler = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		defer GinkgoRecover()
		defer handlerWg.Done()
		conn, err := new(websocket.Upgrader).Upgrade(writer, req, nil)
		conns <- conn
		Expect(err).ToNot(HaveOccurred())

		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			Expect(msgType).To(Equal(websocket.BinaryMessage))

			dopplerMessages <- decodeMessage(data)
		}
	})

	BeforeEach(func() {
		dopplerMessages = make(chan *messages.Doppler, 100)
		conns = make(chan *websocket.Conn, 100)

		mockServer = httptest.NewServer(handler)

		mockDataSource = newMockDataSource()
		mockWaitReporter = newMockWaitReporter()
		mockKvStore = newMockKvStore()

		reg = registry.New(mockKvStore)
		cache = emitter.NewCache(reg)
		dataSourceReader = emitter.NewDataSourceReader(mockDataSource, mockWaitReporter, cache)
	})

	AfterEach(func() {
		close(conns)
		for conn := range conns {
			Expect(conn.Close()).To(Succeed())
		}

		mockServer.CloseClientConnections()
		mockServer.Close()
		handlerWg.Wait()
	})

	Describe("post sharding", func() {
		var (
			callback   func(ID, URL, key string, muxId uint64, add bool)
			expectedId string
		)

		var fetchKvStoreCallback = func() {
			Eventually(mockKvStore.ListenForInput.Callback).Should(Receive(&callback))
		}

		var buildData = func(dataIndex int) []byte {
			return []byte(fmt.Sprintf("some-data-%d", dataIndex))
		}

		var writeDataPoint = func(id string, dataIndex int) {
			dataPoint := &messages.DataPoint{
				Id:   proto.String(id),
				Data: buildData(dataIndex),
			}
			mockDataSource.NextOutput.Ret0 <- dataPoint
			mockDataSource.NextOutput.Ret1 <- true
		}

		BeforeEach(func() {
			expectedId = "some-id"
			writeDataPoint(expectedId, 0)

			fetchKvStoreCallback()
		})

		Context("traffic controller has subscribed", func() {
			var (
				expectedMuxId uint64
				expectedKey   string
			)

			BeforeEach(func() {
				expectedMuxId = 101
				expectedKey = "some-key"

				callback(expectedId, convertHttpToWs(mockServer.URL), expectedKey, expectedMuxId, true)
				handlerWg.Add(1)

				writeDataPoint(expectedId, 1)
			})

			It("sends the expected data point mux ID", func() {
				var msg *messages.Doppler
				Eventually(dopplerMessages).Should(Receive(&msg))

				Expect(msg.GetMuxId()).Should(Equal(expectedMuxId))
			})

			It("sends the expected data point ID", func() {
				var msg *messages.Doppler
				Eventually(dopplerMessages).Should(Receive(&msg))

				Expect(msg.GetMessageType()).Should(Equal(messages.Doppler_DataPoint))
				Expect(msg.GetDataPoint().GetId()).Should(Equal(expectedId))
			})

			It("sends the expected data point data", func() {
				var msg *messages.Doppler
				Eventually(dopplerMessages).Should(Receive(&msg))

				Expect(msg.GetMessageType()).Should(Equal(messages.Doppler_DataPoint))
				Expect(msg.GetDataPoint().GetData()).Should(Equal(buildData(1)))
			})

			Context("several data points are sent", func() {
				var (
					count int
				)

				BeforeEach(func() {
					count = 20
					for i := 2; i < count; i++ {
						writeDataPoint(expectedId, i)
					}
				})

				It("reads all the data points", func() {
					for i := 1; i < count; i++ {
						expectedData := []byte(fmt.Sprintf("some-data-%d", i))
						var msg *messages.Doppler
						Eventually(dopplerMessages).Should(Receive(&msg))

						Expect(msg.GetMessageType()).Should(Equal(messages.Doppler_DataPoint))
						Expect(msg.GetDataPoint().GetData()).Should(Equal(expectedData))
					}
				})
			})
		})
	})
})

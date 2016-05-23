package subscriptionwriter_test

import (
	"fmt"
	"hermes/common/pb/messages"
	"hermes/internal/subscriptionwriter"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionWriter", func() {
	var (
		mockServer *httptest.Server
		handlerWg  sync.WaitGroup

		writer *subscriptionwriter.SubscriptionWriter

		expectedMessage *messages.Doppler
		msgs            chan *messages.Doppler
	)

	var decodeMessage = func(data []byte) *messages.Doppler {
		dataPoint := new(messages.Doppler)
		Expect(dataPoint.Unmarshal(data)).To(Succeed())
		return dataPoint
	}

	var convertUrlToWs = func(URL string) string {
		return fmt.Sprintf("ws%s", URL[4:])
	}

	var handler = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		defer handlerWg.Done()
		handlerWg.Add(1)
		conn, err := new(websocket.Upgrader).Upgrade(writer, req, nil)
		Expect(err).ToNot(HaveOccurred())

		msgType, data, err := conn.ReadMessage()
		Expect(err).ToNot(HaveOccurred())
		Expect(msgType).To(Equal(websocket.BinaryMessage))
		msgs <- decodeMessage(data)
	})

	BeforeEach(func() {
		msgs = make(chan *messages.Doppler, 100)
		mockServer = httptest.NewServer(handler)

		var err error
		writer, err = subscriptionwriter.New(convertUrlToWs(mockServer.URL))
		Expect(err).ToNot(HaveOccurred())

		expectedMessage = &messages.Doppler{
			MuxId:       proto.Uint64(99),
			MessageType: messages.Doppler_DataPoint.Enum(),
			DataPoint: &messages.DataPoint{
				Id:   proto.String("some-id"),
				Data: []byte("some-data"),
			},
		}
	})

	AfterEach(func() {
		mockServer.Close()
		handlerWg.Wait()
	})

	Describe("EmitWithId()", func() {
		It("successfully writes to the websocket", func() {
			Expect(writer.EmitMessage(expectedMessage)).To(Succeed())
		})

		It("writes the expected message", func() {
			writer.EmitMessage(expectedMessage)

			Eventually(msgs).Should(Receive(Equal(expectedMessage)))
		})
	})

})

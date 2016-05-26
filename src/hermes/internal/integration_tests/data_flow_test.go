package integration_test

import (
	"fmt"
	"hermes/common/pb/messages"
	"hermes/internal/emitter"
	"hermes/internal/kvstore/consul"
	"hermes/internal/registry"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("DataFlow", func() {
	var (
		mockServer *httptest.Server
		handlerWg  *sync.WaitGroup

		mockDataSource   *mockDataSource
		mockWaitReporter *mockWaitReporter

		consulStore      *consul.Consul
		dataSourceReader *emitter.DataSourceReader
		cache            *emitter.Cache
		reg              *registry.Registry

		expectedId    string
		expectedMuxId uint64

		dopplerMessages chan *messages.Doppler
		conns           chan *websocket.Conn

		consulSession *gexec.Session
		tmpDir        string
		consulClient  *api.Client
	)

	var _ = AfterSuite(func() {
		consulSession.Kill()
		consulSession.Wait("60s", "200ms")
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
		gexec.CleanupBuildArtifacts()
	})

	var startConsul = func() {
		consulPath, err := gexec.Build("github.com/hashicorp/consul")
		Expect(err).ToNot(HaveOccurred())

		tmpDir, err = ioutil.TempDir("", "consul")
		Expect(err).ToNot(HaveOccurred())

		consulCmd := exec.Command(consulPath, "agent", "-server", "-bootstrap-expect", "1", "-data-dir", tmpDir, "-bind", "127.0.0.1")
		consulSession, err = gexec.Start(consulCmd, nil, nil)
		Expect(err).ToNot(HaveOccurred())
		Consistently(consulSession).ShouldNot(gexec.Exit())

		consulClient, err = api.NewClient(api.DefaultConfig())
		Expect(err).ToNot(HaveOccurred())

		f := func() error {
			_, _, err := consulClient.Catalog().Nodes(nil)
			return err
		}
		Eventually(f, 10).Should(BeNil())
	}

	var decodeMessage = func(data []byte) *messages.Doppler {
		msg := new(messages.Doppler)
		Expect(msg.Unmarshal(data)).To(Succeed())
		return msg
	}

	var convertHttpToWs = func(URL string) string {
		return fmt.Sprintf("ws%s", URL[4:])
	}

	var buildHandler = func(handlerWg *sync.WaitGroup, conns chan *websocket.Conn, dopplerMessages chan *messages.Doppler) http.HandlerFunc {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
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
	}

	var setupIds = func() {
		expectedId = fmt.Sprintf("some-id-%d", rand.Int63())
		expectedMuxId = 101

	}

	BeforeSuite(func() {
		dopplerMessages = make(chan *messages.Doppler, 100)
		conns = make(chan *websocket.Conn, 100)
		handlerWg = new(sync.WaitGroup)
		startConsul()

		mockServer = httptest.NewServer(buildHandler(handlerWg, conns, dopplerMessages))

		mockDataSource = newMockDataSource()
		mockWaitReporter = newMockWaitReporter()

		consulStore = consul.New(convertHttpToWs(mockServer.URL))
		reg = registry.New(consulStore)
		cache = emitter.NewCache(reg)
		dataSourceReader = emitter.NewDataSourceReader(mockDataSource, mockWaitReporter, cache)

		setupIds()

		consulStore.Subscribe(expectedId, expectedMuxId)
		handlerWg.Add(1)

		Eventually(conns, 5).ShouldNot(BeEmpty())
	})

	Describe("post sharding", func() {
		var ()

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

		Context("traffic controller has subscribed", func() {
			var ()

			BeforeEach(func() {
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

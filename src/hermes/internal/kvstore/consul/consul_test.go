package consul_test

import (
	"hermes/internal/kvstore/consul"
	"math/rand"
	"strconv"

	. "github.com/apoydence/eachers"
	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul", func() {
	var (
		clientAddr string
		client     *consul.Consul

		ids    chan string = make(chan string, 100)
		URLs   chan string = make(chan string, 100)
		keys   chan string = make(chan string, 100)
		muxIds chan uint64 = make(chan uint64, 100)
		adds   chan bool   = make(chan bool, 100)
	)

	var callback = func(id, URL, key string, muxId uint64, add bool) {
		ids <- id
		URLs <- URL
		keys <- key
		muxIds <- muxId
		adds <- add
	}

	BeforeEach(func() {
		clientAddr = "some-addr"

		client = consul.New(clientAddr)
	})

	AfterEach(func() {
		_, err := consulClient.KV().DeleteTree(consul.SubscriptionPrefix, nil)
		Expect(err).ToNot(HaveOccurred())

		for len(URLs) > 0 {
			<-URLs
		}
		for len(muxIds) > 0 {
			<-muxIds
		}
		for len(adds) > 0 {
			<-adds
		}
	})

	Describe("Subscribe() & ListenFor()", func() {
		var (
			expectedId1    string
			expectedMuxId1 uint64
			expectedMuxId2 uint64
		)

		BeforeEach(func() {
			expectedId1 = "some-id-1" + strconv.Itoa(int(rand.Int63()))
			expectedMuxId1 = 99
			expectedMuxId2 = 101
		})

		Context("pre-existing announcements", func() {
			BeforeEach(func() {
				client.Subscribe(expectedId1, expectedMuxId1)
				client.Subscribe(expectedId1, expectedMuxId2)

				client.ListenFor(callback)
			})

			It("calls ListenFor callback for each entry", func() {
				Eventually(URLs).Should(HaveLen(2))
			})

			It("sends the expected IDs", func() {
				Eventually(ids).Should(EqualEach(expectedId1, expectedId1))
			})

			It("sends the expected URLs", func() {
				Eventually(URLs).Should(EqualEach(clientAddr, clientAddr))
			})

			It("sends the expected muxIds", func() {
				Eventually(muxIds).Should(BeEquivalentToEach(99, 101))
			})

			It("sends the expected adds", func() {
				Eventually(adds).Should(EqualEach(true, true))
			})
		})
	})

	Describe("Remove()", func() {
		var (
			expectedId1    string
			expectedMuxId1 uint64
			key            string
		)

		var fetchClient = func() *api.Client {
			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).ToNot(HaveOccurred())
			return client
		}

		var fetchCount = func() int {
			pairs, _, err := fetchClient().KV().List(consul.SubscriptionPrefix, nil)
			Expect(err).ToNot(HaveOccurred())
			return len(pairs)
		}

		var drainKeys = func() {
			for len(keys) > 0 {
				<-keys
			}
		}

		BeforeEach(func() {
			expectedId1 = "some-id-1" + strconv.Itoa(int(rand.Int63()))
			expectedMuxId1 = 99
			drainKeys()

			client.Subscribe(expectedId1, expectedMuxId1)
			Eventually(keys).Should(Receive(&key))
		})

		It("removes the entry", func() {
			client.Remove(key)

			Eventually(fetchCount).Should(Equal(0))
		})
	})
})

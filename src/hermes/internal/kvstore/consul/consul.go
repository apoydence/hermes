package consul

import (
	"fmt"
	"hermes/common/pb/kvstore"
	"math/rand"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
)

const (
	SubscriptionPrefix = "Subscription"
)

type Consul struct {
	clientAddr string
	client     *api.Client
	kv         *api.KV
}

func New(clientAddr string) *Consul {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(fmt.Sprintf("Unable to create client: %s", err))
	}

	return &Consul{
		clientAddr: clientAddr,
		client:     client,
		kv:         client.KV(),
	}
}

func (c *Consul) Subscribe(id string, muxId uint64) {
	pair := &api.KVPair{
		Key:   buildUniqueKey(id),
		Value: c.buildValue(id, muxId),
	}

	if _, err := c.kv.Put(pair, nil); err != nil {
		panic(fmt.Sprintf("unable to subscribe to %s (muxId=%d): %s", id, muxId, err))
	}
}

func (c *Consul) ListenFor(id string, callback func(id, URL, key string, muxId uint64, add bool)) {
	go c.listenFor(id, callback)
}

func (c *Consul) Remove(key string) {
	_, err := c.kv.Delete(key, nil)
	if err != nil {
		fmt.Println("Unable to delete key", err)
	}
}

func (c *Consul) listenFor(id string, callback func(id, URL, key string, muxId uint64, add bool)) {
	var options *api.QueryOptions

	for {
		pairs, meta, err := c.kv.List(SubscriptionPrefix, options)
		if err != nil {
			fmt.Println("Unable to list keys", err)
			return
		}

		if options == nil {
			options = new(api.QueryOptions)
		}
		options.WaitIndex = meta.LastIndex

		for _, pair := range pairs {
			sub := c.readValue(pair)
			if sub == nil {
				continue
			}
			callback(sub.GetId(), sub.GetUrl(), pair.Key, sub.GetMuxId(), true)
		}
	}
}

func (c *Consul) buildValue(id string, muxId uint64) []byte {
	sub := &kvstore.Subscription{
		Id:    proto.String(id),
		MuxId: proto.Uint64(muxId),
		Url:   proto.String(c.clientAddr),
	}
	data, err := sub.Marshal()
	if err != nil {
		panic(err)
	}

	return data
}

func (c *Consul) readValue(pair *api.KVPair) *kvstore.Subscription {
	sub := new(kvstore.Subscription)
	if err := sub.Unmarshal(pair.Value); err != nil {
		return nil
	}
	return sub
}

func buildUniqueKey(id string) string {
	return fmt.Sprintf("%s-%s-%020d", SubscriptionPrefix, id, uint64(rand.Int63()))
}

func stripSubscriptionMeta(key string) string {
	return key[len(SubscriptionPrefix)+1 : len(key)-21]
}

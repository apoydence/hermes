package emitter

import "hermes/internal/tcwriter"

type KvStore interface {
	ListenFor(id string, callback func(URL string, muxId uint64))
}

type Generator struct {
	kvstore KvStore
}

func NewGenerator(kvstore KvStore) *Generator {
	return &Generator{
		kvstore: kvstore,
	}
}

func (g *Generator) Fetch(id string) Emitter {
	ll := NewLinkedList()
	g.kvstore.ListenFor(id, func(URL string, muxId uint64) {
		connWriter, err := tcwriter.New(URL)
		if err != nil {
			panic(err)
		}
		ll.Append(tcwriter.NewEmitter(muxId, connWriter))
	})

	return NewSubscriptionListReader(ll)
}

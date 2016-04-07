package routing

import "hermes/internal/tcwriter"

type KvStore interface {
	ListenFor(id string, callback func(URL string, muxId uint64))
}

type EmitterGenerator struct {
	kvstore KvStore
}

func NewEmitterGenerator(kvstore KvStore) *EmitterGenerator {
	return &EmitterGenerator{
		kvstore: kvstore,
	}
}

func (g *EmitterGenerator) Fetch(id string) Emitter {
	ll := NewLinkedList()
	g.kvstore.ListenFor(id, func(URL string, muxId uint64) {
		connWriter, err := tcwriter.New(URL)
		if err != nil {
			panic(err)
		}
		ll.Append(tcwriter.NewEmitter(muxId, connWriter))
	})

	return NewEventEmitter(ll)
}

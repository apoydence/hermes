package registry

import (
	"hermes/internal/datastructures"
	"hermes/internal/subscriptionwriter"
	"sync"
	"unsafe"
)

type KvStore interface {
	ListenFor(callback func(ID, URL, key string, muxID uint64, add bool))
	Remove(key string)
}

type Registry struct {
	lists   map[string]*datastructures.LinkedList
	kvstore KvStore
	lock    sync.Mutex
}

func New(kvstore KvStore) *Registry {
	r := &Registry{
		lists:   make(map[string]*datastructures.LinkedList),
		kvstore: kvstore,
	}
	kvstore.ListenFor(r.listen)
	return r
}

func (r *Registry) GetList(ID string) *datastructures.LinkedList {
	r.lock.Lock()
	defer r.lock.Unlock()
	if list, ok := r.lists[ID]; ok {
		return list
	}

	list := datastructures.NewLinkedList()
	r.lists[ID] = list
	return list
}

func (r *Registry) listen(ID, URL, key string, muxID uint64, add bool) {
	list := r.GetList(ID)
	writer, err := subscriptionwriter.New(URL)
	if err != nil {
		r.kvstore.Remove(key)
	}
	emitter := subscriptionwriter.NewEmitter(muxID, writer)
	list.Append(unsafe.Pointer(emitter))
}

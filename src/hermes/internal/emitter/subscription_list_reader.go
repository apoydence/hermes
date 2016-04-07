package emitter

import (
	"hermes/common/pb/messages"
	"unsafe"
)

type SenderStore interface {
	Traverse(callback func(sender unsafe.Pointer))
}

type SubscriptionListReader struct {
	senderStore SenderStore
}

var NewSubscriptionListReader = func(senderStore SenderStore) *SubscriptionListReader {
	return &SubscriptionListReader{
		senderStore: senderStore,
	}
}

func (r *SubscriptionListReader) Emit(data *messages.DataPoint) error {
	r.senderStore.Traverse(func(sender unsafe.Pointer) {
		emitter := *(*Emitter)(sender)
		emitter.Emit(data)
	})
	return nil
}

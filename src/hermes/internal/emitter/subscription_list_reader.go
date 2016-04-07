package emitter

import "hermes/common/pb/messages"

type SenderStore interface {
	Traverse(callback func(sender Emitter))
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
	r.senderStore.Traverse(func(sender Emitter) {
		sender.Emit(data)
	})
	return nil
}

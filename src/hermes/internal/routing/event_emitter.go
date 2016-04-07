package routing

import "hermes/common/pb/messages"

type SenderStore interface {
	Traverse(callback func(sender Emitter))
}

type EventEmitter struct {
	senderStore SenderStore
}

var NewEventEmitter = func(senderStore SenderStore) *EventEmitter {
	return &EventEmitter{
		senderStore: senderStore,
	}
}

func (e *EventEmitter) Emit(data *messages.DataPoint) error {
	e.senderStore.Traverse(func(sender Emitter) {
		sender.Emit(data)
	})
	return nil
}

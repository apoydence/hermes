package tcwriter

import (
	"hermes/common/pb/messages"

	"github.com/gogo/protobuf/proto"
)

type MessageEmitter interface {
	EmitMessage(msg *messages.Doppler) error
}

type Emitter struct {
	muxId      uint64
	msgEmitter MessageEmitter
}

var NewEmitter = func(muxId uint64, msgEmitter MessageEmitter) *Emitter {
	return &Emitter{
		muxId:      muxId,
		msgEmitter: msgEmitter,
	}
}

func (e *Emitter) Emit(dataPoint *messages.DataPoint) error {
	return e.msgEmitter.EmitMessage(e.buildMessage(dataPoint))
}

func (e *Emitter) buildMessage(dataPoint *messages.DataPoint) *messages.Doppler {
	return &messages.Doppler{
		MuxId:       proto.Uint64(e.muxId),
		MessageType: messages.Doppler_DataPoint.Enum(),
		DataPoint:   dataPoint,
	}
}

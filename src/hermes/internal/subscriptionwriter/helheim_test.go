package subscriptionwriter_test

import "hermes/common/pb/messages"

type mockMessageEmitter struct {
	EmitMessageCalled chan bool
	EmitMessageInput  struct {
		Msg chan *messages.Doppler
	}
	EmitMessageOutput struct {
		Ret0 chan error
	}
}

func newMockMessageEmitter() *mockMessageEmitter {
	m := &mockMessageEmitter{}
	m.EmitMessageCalled = make(chan bool, 100)
	m.EmitMessageInput.Msg = make(chan *messages.Doppler, 100)
	m.EmitMessageOutput.Ret0 = make(chan error, 100)
	return m
}
func (m *mockMessageEmitter) EmitMessage(msg *messages.Doppler) error {
	m.EmitMessageCalled <- true
	m.EmitMessageInput.Msg <- msg
	return <-m.EmitMessageOutput.Ret0
}

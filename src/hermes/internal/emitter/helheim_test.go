package emitter_test

import (
	"hermes/common/pb/messages"
	"hermes/internal/emitter"
)

type mockKvStore struct {
	ListenForCalled chan bool
	ListenForInput  struct {
		Id       chan string
		Callback chan func(URL string, muxId uint64)
	}
}

func newMockKvStore() *mockKvStore {
	m := &mockKvStore{}
	m.ListenForCalled = make(chan bool, 100)
	m.ListenForInput.Id = make(chan string, 100)
	m.ListenForInput.Callback = make(chan func(URL string, muxId uint64), 100)
	return m
}
func (m *mockKvStore) ListenFor(id string, callback func(URL string, muxId uint64)) {
	m.ListenForCalled <- true
	m.ListenForInput.Id <- id
	m.ListenForInput.Callback <- callback
}

type mockSenderStore struct {
	TraverseCalled chan bool
	TraverseInput  struct {
		Callback chan func(sender emitter.Emitter)
	}
}

func newMockSenderStore() *mockSenderStore {
	m := &mockSenderStore{}
	m.TraverseCalled = make(chan bool, 100)
	m.TraverseInput.Callback = make(chan func(sender emitter.Emitter), 100)
	return m
}
func (m *mockSenderStore) Traverse(callback func(sender emitter.Emitter)) {
	m.TraverseCalled <- true
	m.TraverseInput.Callback <- callback
}

type mockEmitterFetcher struct {
	FetchCalled chan bool
	FetchInput  struct {
		Id chan string
	}
	FetchOutput struct {
		Ret0 chan emitter.Emitter
	}
}

func newMockEmitterFetcher() *mockEmitterFetcher {
	m := &mockEmitterFetcher{}
	m.FetchCalled = make(chan bool, 100)
	m.FetchInput.Id = make(chan string, 100)
	m.FetchOutput.Ret0 = make(chan emitter.Emitter, 100)
	return m
}
func (m *mockEmitterFetcher) Fetch(id string) emitter.Emitter {
	m.FetchCalled <- true
	m.FetchInput.Id <- id
	return <-m.FetchOutput.Ret0
}

type mockDataSource struct {
	NextCalled chan bool
	NextOutput struct {
		Ret0 chan *messages.DataPoint
		Ret1 chan bool
	}
}

func newMockDataSource() *mockDataSource {
	m := &mockDataSource{}
	m.NextCalled = make(chan bool, 100)
	m.NextOutput.Ret0 = make(chan *messages.DataPoint, 100)
	m.NextOutput.Ret1 = make(chan bool, 100)
	return m
}
func (m *mockDataSource) Next() (*messages.DataPoint, bool) {
	m.NextCalled <- true
	return <-m.NextOutput.Ret0, <-m.NextOutput.Ret1
}

type mockWaitReporter struct {
	WaitingCalled chan bool
}

func newMockWaitReporter() *mockWaitReporter {
	m := &mockWaitReporter{}
	m.WaitingCalled = make(chan bool, 100)
	return m
}
func (m *mockWaitReporter) Waiting() {
	m.WaitingCalled <- true
}

type mockEmitter struct {
	EmitCalled chan bool
	EmitInput  struct {
		Data chan *messages.DataPoint
	}
	EmitOutput struct {
		Ret0 chan error
	}
}

func newMockEmitter() *mockEmitter {
	m := &mockEmitter{}
	m.EmitCalled = make(chan bool, 100)
	m.EmitInput.Data = make(chan *messages.DataPoint, 100)
	m.EmitOutput.Ret0 = make(chan error, 100)
	return m
}
func (m *mockEmitter) Emit(data *messages.DataPoint) error {
	m.EmitCalled <- true
	m.EmitInput.Data <- data
	return <-m.EmitOutput.Ret0
}

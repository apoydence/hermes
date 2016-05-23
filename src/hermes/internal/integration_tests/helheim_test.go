package integration_test

import "hermes/common/pb/messages"

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

type mockKvStore struct {
	ListenForCalled chan bool
	ListenForInput  struct {
		Callback chan func(ID, URL, key string, muxID uint64, add bool)
	}
	RemoveCalled chan bool
	RemoveInput  struct {
		Key chan string
	}
}

func newMockKvStore() *mockKvStore {
	m := &mockKvStore{}
	m.ListenForCalled = make(chan bool, 100)
	m.ListenForInput.Callback = make(chan func(ID, URL, key string, muxID uint64, add bool), 100)
	m.RemoveCalled = make(chan bool, 100)
	m.RemoveInput.Key = make(chan string, 100)
	return m
}
func (m *mockKvStore) ListenFor(callback func(ID, URL, key string, muxID uint64, add bool)) {
	m.ListenForCalled <- true
	m.ListenForInput.Callback <- callback
}
func (m *mockKvStore) Remove(key string) {
	m.RemoveCalled <- true
	m.RemoveInput.Key <- key
}

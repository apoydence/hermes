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

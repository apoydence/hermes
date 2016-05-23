// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package registry_test

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
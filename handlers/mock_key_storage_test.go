package handlers_test

import (
	"errors"
	"github.com/poy/hermes/handlers"
	"io"
)

type mockKeyStorage map[string]*handlers.ReadLocker

func NewMockKeyStorage() mockKeyStorage {
	return make(map[string]*handlers.ReadLocker)
}

func (m mockKeyStorage) Add(key string, data io.Reader) error {
	if _, ok := m[key]; ok {
		return errors.New("Key already present")
	}
	m[key] = &handlers.ReadLocker{
		Reader: data,
		Lock:   handlers.NewLocker(),
	}
	return nil
}

func (m mockKeyStorage) Fetch(key string) *handlers.ReadLocker {
	if c, ok := m[key]; ok {
		return c
	}
	return nil
}

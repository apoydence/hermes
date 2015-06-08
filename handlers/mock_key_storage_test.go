package handlers_test

import (
	"errors"
	"io"
)

type mockKeyStorage map[string]io.Reader

func NewMockKeyStorage() mockKeyStorage {
	return make(map[string]io.Reader)
}

func (m mockKeyStorage) Add(key string, data io.Reader) error {
	if _, ok := m[key]; ok {
		return errors.New("Key already present")
	}
	m[key] = data
	return nil
}

func (m mockKeyStorage) Fetch(key string) io.Reader {
	if c, ok := m[key]; ok {
		return c
	}
	return nil
}

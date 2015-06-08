package handlers

import (
	"fmt"
	"io"
	"sync"
)

type KeyStorer struct {
	rw   *sync.RWMutex
	keys map[string]io.Reader
}

func NewKeyStorer() *KeyStorer {
	return &KeyStorer{
		rw:   &sync.RWMutex{},
		keys: make(map[string]io.Reader),
	}
}

func (k *KeyStorer) Add(key string, reader io.Reader) error {
	k.rw.Lock()
	defer k.rw.Unlock()
	if _, ok := k.keys[key]; ok {
		return fmt.Errorf("The key %s is already in use", key)
	}
	k.keys[key] = reader
	return nil
}

func (k *KeyStorer) Fetch(key string) io.Reader {
	k.rw.RLock()
	defer k.rw.RUnlock()
	if reader, ok := k.keys[key]; ok {
		return reader
	}
	return nil
}

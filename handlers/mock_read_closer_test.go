package handlers_test

import (
	"bytes"
	"io"
)

type mockReadCloser struct {
	buffer    *bytes.Buffer
	isClosed  bool
	doneAfter int
}

func NewMockReadCloser(buf *bytes.Buffer, doneAfter int) *mockReadCloser {
	return &mockReadCloser{
		buffer:    buf,
		doneAfter: doneAfter,
	}
}

func (m *mockReadCloser) Read(buffer []byte) (int, error) {
	if m.isClosed {
		panic("Already closed")
	}
	if m.doneAfter <= 0 {
		return 0, io.EOF
	}
	m.doneAfter--
	return m.buffer.Read(buffer)
}

func (m *mockReadCloser) Close() error {
	m.isClosed = true
	return nil
}

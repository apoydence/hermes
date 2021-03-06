package handlers

import (
	"io"
	"time"
)

type channelReader struct {
	done    bool
	channel chan byte
	reader  io.ReadCloser
}

func NewChannelReader(reader io.ReadCloser) *channelReader {
	c := &channelReader{
		channel: make(chan byte, 0),
		reader:  reader,
	}
	return c
}

func (c *channelReader) Read(buffer []byte) (int, error) {
	if c.done {
		return 0, io.EOF
	}
	i := 0
	for ; i < len(buffer); i++ {
		d, ok := <-c.channel
		if !ok {
			c.done = true
			return i, nil
		}
		buffer[i] = d
	}
	return i, nil
}

func (c *channelReader) Run(timeout time.Duration) {
	defer c.reader.Close()
	defer close(c.channel)
	buffer := make([]byte, 1024)
	timer := time.NewTimer(timeout)
	for {
		i, err := c.reader.Read(buffer)
		if err != nil && (err != io.EOF || i == 0) {
			return
		}

		for j := 0; j < i; j++ {
			select {
			case <-timer.C:
				return
			case c.channel <- buffer[j]:
				timer.Reset(timeout)
			}
		}
	}
}

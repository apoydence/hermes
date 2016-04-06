package tcwriter

import (
	"hermes/common/pb/messages"

	"github.com/gorilla/websocket"
)

type TcWriter struct {
	URL  string
	conn *websocket.Conn
}

var New = func(URL string) (*TcWriter, error) {
	writer := &TcWriter{
		URL: URL,
	}

	return writer, writer.connect()
}

func (w *TcWriter) EmitMessage(msg *messages.Doppler) error {
	data, err := msg.Marshal()
	if err != nil {
		panic(err)
	}

	return w.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (w *TcWriter) Close() error {
	return w.conn.Close()
}

func (w *TcWriter) connect() error {
	var err error
	w.conn, _, err = websocket.DefaultDialer.Dial(w.URL, nil)
	return err
}

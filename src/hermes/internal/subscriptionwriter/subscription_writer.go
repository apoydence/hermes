package subscriptionwriter

import (
	"hermes/common/pb/messages"

	"github.com/gorilla/websocket"
)

type SubscriptionWriter struct {
	URL  string
	conn *websocket.Conn
}

var New = func(URL string) (*SubscriptionWriter, error) {
	writer := &SubscriptionWriter{
		URL: URL,
	}

	return writer, writer.connect()
}

func (w *SubscriptionWriter) EmitMessage(msg *messages.Doppler) error {
	data, err := msg.Marshal()
	if err != nil {
		panic(err)
	}

	return w.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (w *SubscriptionWriter) Close() error {
	return w.conn.Close()
}

func (w *SubscriptionWriter) connect() error {
	var err error
	w.conn, _, err = websocket.DefaultDialer.Dial(w.URL, nil)
	return err
}

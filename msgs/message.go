package msgs

import (
	"time"
)

type Message struct {
	id   int64
	from string
	data []byte
}

func NewMessage(from string, data []byte) *Message {
	mes := &Message{
		id:   time.Now().UnixNano(),
		from: from,
		data: data,
	}
	return mes
}

func (m *Message) ID() int64 {
	return m.id
}

func (m *Message) From() string {
	return m.from
}

func (m *Message) Data() []byte {
	return m.data
}

// Package msgs implements entities Message and MessageQueue
// Message keeps some metadata and also payload
package msgs

import (
	"time"
)

type Message struct {
	id   int64
	from string
	data []byte
}

// NewMessage creates message
func NewMessage(from string, data []byte) *Message {
	mes := &Message{
		id:   time.Now().UnixNano(),
		from: from,
		data: data,
	}
	return mes
}

// ID getter id field
func (m *Message) ID() int64 {
	return m.id
}

// From getter from field
func (m *Message) From() string {
	return m.from
}

// Data getter data field
func (m *Message) Data() []byte {
	return m.data
}

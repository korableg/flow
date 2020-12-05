package repo

import "io"

type DB interface {
	NewNodeRepository(name ...string) NodeDB
	NewHubRepository() HubDB
	NewMQRepository(name string) MQDB
}

type NodeDB interface {
	Parent() DB
	Store(*Node) error
	All() ([]*Node, error)
	Delete(string) error
	DeleteDB() error
	io.Closer
}

type HubDB interface {
	Parent() DB
	Store(*Hub) error
	All() ([]*Hub, error)
	Delete(string) error
	io.Closer
}

type MQDB interface {
	Parent() DB
	Store(*Message) error
	StoreAll(MessageSlice) error
	All() (MessageSlice, error)
	Delete(int64) error
	DeleteDB() error
	io.Closer
}

type Node struct {
	Name    string
	Careful bool
}

type Hub struct {
	Name string
}

type Message struct {
	ID   int64
	From string
	Data []byte
}

type MessageSlice []*Message

func (ms MessageSlice) Len() int {
	return len(ms)
}

func (ms MessageSlice) Less(i, j int) bool {
	return ms[i].ID < ms[j].ID
}

func (ms MessageSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

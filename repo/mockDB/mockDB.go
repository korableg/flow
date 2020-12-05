package MockDB

import (
	"github.com/korableg/mini-gin/flow/repo"
)

type MockDB struct{}

func (t *MockDB) NewNodeRepository(hubName ...string) repo.NodeDB {
	n := new(mockNodeRepository)
	n.parent = t
	return n
}

func (t *MockDB) NewHubRepository() repo.HubDB {
	n := new(mockHubRepository)
	n.parent = t
	return n
}

func (t *MockDB) NewMQRepository(name string) repo.MQDB {
	n := new(mockMQRepository)
	n.parent = t
	return n
}

type mockHubRepository struct {
	parent repo.DB
}

func (hr *mockHubRepository) Parent() repo.DB {
	return hr.parent
}

func (hr *mockHubRepository) Store(hub *repo.Hub) error {
	return nil
}

func (hr *mockHubRepository) All() ([]*repo.Hub, error) {
	hubs := make([]*repo.Hub, 0, 20)
	hubs = append(hubs, &repo.Hub{Name: "mockHub"})
	return hubs, nil
}

func (hr *mockHubRepository) Delete(key string) error {
	return nil
}

func (hr *mockHubRepository) Close() error {
	return nil
}

type mockNodeRepository struct {
	parent repo.DB
}

func (nr *mockNodeRepository) Parent() repo.DB {
	return nr.parent
}

func (nr *mockNodeRepository) Store(node *repo.Node) error {
	return nil
}

func (nr *mockNodeRepository) All() ([]*repo.Node, error) {
	nodes := make([]*repo.Node, 0, 20)
	nodes = append(nodes, &repo.Node{Name: "mockNode"})
	return nodes, nil
}

func (nr *mockNodeRepository) Delete(key string) error {
	return nil
}

func (nr *mockNodeRepository) DeleteDB() error {
	return nil
}

func (nr *mockNodeRepository) Close() error {
	return nil
}

type mockMQRepository struct {
	parent repo.DB
}

func (mqr *mockMQRepository) Parent() repo.DB {
	return mqr.parent
}

func (mqr *mockMQRepository) Store(mes *repo.Message) error {
	return nil
}

func (mqr *mockMQRepository) StoreAll(mes repo.MessageSlice) error {
	return nil
}

func (mqr *mockMQRepository) All() (repo.MessageSlice, error) {
	messages := make(repo.MessageSlice, 0, 20)
	messages = append(messages, &repo.Message{ID: 3, From: "testfrom1", Data: nil})
	messages = append(messages, &repo.Message{ID: 1, From: "testfrom1", Data: nil})
	messages = append(messages, &repo.Message{ID: 15, From: "testfrom1", Data: nil})
	messages = append(messages, &repo.Message{ID: 47, From: "testfrom1", Data: nil})
	return messages, nil
}

func (mqr *mockMQRepository) Delete(id int64) error {
	return nil
}

func (mqr *mockMQRepository) DeleteDB() error {
	return nil
}

func (mqr *mockMQRepository) Close() error {
	return nil
}

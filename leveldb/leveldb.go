package leveldb

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/korableg/mini-gin/flow/repo"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path/filepath"
	"strings"
)

type GoLevelDB struct {
	path string
}

func New(path string) *GoLevelDB {
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	db := GoLevelDB{
		path: path,
	}
	return &db
}

func (f *GoLevelDB) NewNodeRepository(name ...string) repo.NodeDB {

	path := "nodes"
	if name != nil && len(name) > 0 {
		path = fmt.Sprintf("%s%c%s", "nodesinhubs", filepath.Separator, name[0])
	}
	dbPath := fmt.Sprintf("%s%c%s%c%s", f.path, filepath.Separator, "db", filepath.Separator, path)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	n := NewNodeRepository(db, dbPath, f)
	return n
}

func (f *GoLevelDB) NewHubRepository() repo.HubDB {

	dbPath := fmt.Sprintf("%s%c%s%c%s", f.path, filepath.Separator, "db", filepath.Separator, "hubs")

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	n := NewHubRepository(db, f)
	return n
}

func (f *GoLevelDB) NewMQRepository(name string) repo.MQDB {

	dbPath := fmt.Sprintf("%s%c%s%c%s%c%s",
		f.path, filepath.Separator, "db", filepath.Separator, "messages", filepath.Separator, name)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	n := NewMQRepository(db, dbPath, f)
	return n
}

type HubRepository struct {
	db     *leveldb.DB
	parent repo.DB
}

func NewHubRepository(db *leveldb.DB, parent repo.DB) *HubRepository {
	hr := HubRepository{db: db, parent: parent}
	return &hr
}

func (hr *HubRepository) Parent() repo.DB {
	return hr.parent
}

func (hr *HubRepository) Store(hub *repo.Hub) error {

	dataJSON, err := json.Marshal(hub)
	if err != nil {
		return err
	}
	err = hr.db.Put([]byte(hub.Name), dataJSON, nil)

	return err

}

func (hr *HubRepository) All() ([]*repo.Hub, error) {

	hubs := make([]*repo.Hub, 0, 20)

	iterator := hr.db.NewIterator(nil, nil)
	for iterator.Next() {
		hub := new(repo.Hub)
		err := json.Unmarshal(iterator.Value(), hub)
		if err != nil {
			return nil, err
		}
		hubs = append(hubs, hub)
	}
	iterator.Release()

	if err := iterator.Error(); err != nil {
		return nil, err
	}

	return hubs, nil

}

func (hr *HubRepository) Delete(key string) error {
	return hr.db.Delete([]byte(key), nil)
}

func (hr *HubRepository) Close() error {
	return hr.db.Close()
}

type NodeRepository struct {
	db     *leveldb.DB
	path   string
	parent repo.DB
}

func NewNodeRepository(db *leveldb.DB, path string, parent repo.DB) *NodeRepository {
	nr := NodeRepository{db: db, path: path, parent: parent}
	return &nr
}

func (nr *NodeRepository) Parent() repo.DB {
	return nr.parent
}

func (nr *NodeRepository) Store(node *repo.Node) error {

	dataJSON, err := json.Marshal(node)
	if err != nil {
		return err
	}
	err = nr.db.Put([]byte(node.Name), dataJSON, nil)

	return err

}

func (nr *NodeRepository) All() ([]*repo.Node, error) {

	nodes := make([]*repo.Node, 0, 20)

	iterator := nr.db.NewIterator(nil, nil)
	for iterator.Next() {
		node := new(repo.Node)
		err := json.Unmarshal(iterator.Value(), node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	iterator.Release()

	if err := iterator.Error(); err != nil {
		return nil, err
	}

	return nodes, nil

}

func (nr *NodeRepository) Delete(key string) error {
	return nr.db.Delete([]byte(key), nil)
}

func (nr *NodeRepository) DeleteDB() (err error) {
	if err = nr.Close(); err == nil {
		err = os.RemoveAll(nr.path)
	}
	return
}

func (nr *NodeRepository) Close() error {
	return nr.db.Close()
}

type MQRepository struct {
	db     *leveldb.DB
	path   string
	parent repo.DB
}

func NewMQRepository(db *leveldb.DB, path string, parent repo.DB) *MQRepository {
	r := MQRepository{db: db, path: path, parent: parent}
	return &r
}

func (t *MQRepository) Parent() repo.DB {
	return t.parent
}

func (t *MQRepository) Store(message *repo.Message) error {

	dataJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = t.db.Put(t.int64tobyteslice(message.ID), dataJSON, nil)

	return err

}

func (t *MQRepository) StoreAll(messages repo.MessageSlice) error {

	batch := new(leveldb.Batch)
	for _, message := range messages {
		dataJSON, err := json.Marshal(message)
		if err != nil {
			return err
		}
		batch.Put(t.int64tobyteslice(message.ID), dataJSON)
	}

	return t.db.Write(batch, nil)

}

func (t *MQRepository) All() (repo.MessageSlice, error) {

	messages := make(repo.MessageSlice, 0)

	iterator := t.db.NewIterator(nil, nil)
	for iterator.Next() {
		message := new(repo.Message)
		err := json.Unmarshal(iterator.Value(), message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	iterator.Release()

	if err := iterator.Error(); err != nil {
		return nil, err
	}

	return messages, nil

}

func (t *MQRepository) Delete(ID int64) error {
	return t.db.Delete(t.int64tobyteslice(ID), nil)
}

func (t *MQRepository) DeleteDB() (err error) {
	if err = t.Close(); err == nil {
		err = os.RemoveAll(t.path)
	}
	return
}

func (t *MQRepository) Close() error {
	return t.db.Close()
}

func (t *MQRepository) int64tobyteslice(value int64) []byte {
	byteID := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteID, uint64(value))
	return byteID
}

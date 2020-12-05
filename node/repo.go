package node

import (
	"github.com/korableg/mini-gin/flow/repo"
	"sync"
)

type Repository struct {
	nodes *sync.Map
	db    repo.NodeDB
}

func NewNodeRepository(db repo.NodeDB, nodeTmplts ...*Repository) *Repository {

	nr := new(Repository)
	nr.nodes = new(sync.Map)
	nr.db = db

	if nr.db == nil {
		return nr
	}

	var nodesTmplt *Repository
	if nodeTmplts != nil && len(nodeTmplts) > 0 {
		nodesTmplt = nodeTmplts[0]
	}
	nodesRepo, err := nr.db.All()
	if err != nil {
		panic(err)
	}
	for _, nodeRepo := range nodesRepo {

		var n *Node
		if nodesTmplt != nil {
			n, _ = nodesTmplt.Load(nodeRepo.Name)
		}
		if n == nil {
			n, err = New(nodeRepo.Name, nodeRepo.Careful, db.Parent())
			if err != nil {
				panic(err)
			}
		}

		nr.nodes.Store(n.Name(), n)
	}

	return nr
}

func (nr *Repository) Store(node *Node) error {
	if nr.db != nil {
		nodeRepo := new(repo.Node)
		nodeRepo.Name = node.Name()
		nodeRepo.Careful = node.IsCareful()
		if err := nr.db.Store(nodeRepo); err != nil {
			return err
		}
	}
	nr.nodes.Store(node.Name(), node)
	return nil
}

func (nr *Repository) Load(name string) (*Node, bool) {
	if n, ok := nr.nodes.Load(name); ok {
		return n.(*Node), ok
	}
	return nil, false
}

func (nr *Repository) Slice() []*Node {
	nodes := make([]*Node, 0, 20)
	f := func(value *Node) error { nodes = append(nodes, value); return nil }
	_ = nr.Range(f)
	return nodes
}

func (nr *Repository) Range(f func(node *Node) error) error {
	var err error
	rangeFunc := func(key, value interface{}) bool {
		err = f(value.(*Node))
		if err != nil {
			return false
		}
		return true
	}
	nr.nodes.Range(rangeFunc)
	return err
}

func (nr *Repository) Delete(name string) error {
	if nr.db != nil {
		if err := nr.db.Delete(name); err != nil {
			return err
		}
	}
	nr.nodes.Delete(name)
	return nil
}

func (nr *Repository) DeleteDB() error {
	if nr.db == nil {
		return nil
	}
	err := nr.db.DeleteDB()
	if err != nil {
		return err
	}
	nr.db = nil
	return nil
}

func (nr *Repository) Close() (err error) {
	if nr.db == nil {
		return
	}

	wg := new(sync.WaitGroup)

	f := func(n *Node) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = n.messages.Close() // the error isn't handle
		}()
		return nil
	}

	_ = nr.Range(f) // the error isn't handle because it always nil

	err = nr.db.Close()

	wg.Wait()

	return
}

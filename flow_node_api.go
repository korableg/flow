// Package flow implements basic methods to working with framework
package flow

import (
	"github.com/korableg/flow/errs"
	"github.com/korableg/flow/hub"
	"github.com/korableg/flow/node"
)

// NewNode creates new Node by name. If careful is true when every messages will be saved to disk
func (m *Flow) NewNode(name string, careful bool) (*node.Node, error) {

	if nodeExists := m.GetNode(name); nodeExists != nil {
		return nil, errs.ErrNodeIsAlreadyExists
	}

	n, err := node.New(name, careful, m.db)
	if err != nil {
		return nil, err
	}

	if err := m.nodes.Store(n); err != nil {
		return nil, err
	}

	return n, nil

}

// GetNode gets node by name
func (m *Flow) GetNode(name string) (n *node.Node) {
	if value, ok := m.nodes.Load(name); ok {
		n = value
	}
	return
}

// GetAllNodes gets slice of nodes
func (m *Flow) GetAllNodes() []*node.Node {
	return m.nodes.Slice()
}

// DeleteNode removes node by name
func (m *Flow) DeleteNode(name string) error {
	n := m.GetNode(name)
	if n == nil {
		return nil
	}
	err := n.DeleteMessageDB()
	if err != nil {
		return err
	}
	f := func(hub *hub.Hub) error {
		return hub.DeleteNode(n)
	}
	err = m.hubs.Range(f)
	if err != nil {
		return err
	}
	return m.nodes.Delete(name)
}

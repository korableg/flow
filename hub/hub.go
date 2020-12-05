// Package hub implements entities Hub and Repository.
// Hub sends message to all include nodes.
package hub

import (
	"encoding/json"
	"github.com/korableg/flow/cmn"
	"github.com/korableg/flow/errs"
	"github.com/korableg/flow/msgs"
	"github.com/korableg/flow/node"
	"github.com/korableg/flow/repo"
)

type Hub struct {
	name  string
	nodes *node.Repository
}

// New creates new hub
func New(name string, db repo.DB, nodes ...*node.Repository) (h *Hub, err error) {

	if err = checkName(name); err != nil {
		return
	}

	var nodeDB repo.NodeDB
	if db != nil {
		nodeDB = db.NewNodeRepository(name)
	}

	h = new(Hub)
	h.name = name
	h.nodes = node.NewNodeRepository(nodeDB, nodes...)

	return

}

// Name getter
func (h *Hub) Name() string {
	return h.name
}

// AddNode adding node to hub
func (h *Hub) AddNode(n *node.Node) error {
	return h.nodes.Store(n)
}

// DeleteNode deleting node from hub
func (h *Hub) DeleteNode(n *node.Node) error {
	if n == nil {
		return nil
	}
	return h.nodes.Delete(n.Name())
}

// PushMessage sending message in to each node of hub
func (h *Hub) PushMessage(m *msgs.Message) error {
	return h.nodes.Range(func(n *node.Node) error { return n.PushMessage(m) })
}

func (h *Hub) MarshalJSON() ([]byte, error) {

	nodes := h.nodes.Slice()

	hubMap := make(map[string]interface{})
	hubMap["name"] = h.name
	hubMap["nodes"] = nodes

	return json.Marshal(hubMap)

}

// deleteAllNodes deleting all nodes from this hub
func (h *Hub) deleteAllNodes() error {
	nodes := h.nodes.Slice()
	for _, n := range nodes {
		if err := h.DeleteNode(n); err != nil {
			return err
		}
	}
	return nil
}

// deleteNodeDB deleting node db, usually before deleting this hub
func (h *Hub) deleteNodeDB() error {
	err := h.nodes.DeleteDB()
	return err
}

// checkName checking hub name on match some rules
func checkName(name string) error {
	if len(name) == 0 {
		return errs.ErrHubNameIsempty
	}
	if len([]rune(name)) > 100 {
		return errs.ErrHubNameOver100
	}
	if !cmn.NameMatchedPattern(name) {
		return errs.ErrHubNameNotMatchedPattern
	}
	return nil
}

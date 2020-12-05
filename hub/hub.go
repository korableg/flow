package hub

import (
	"encoding/json"
	"github.com/korableg/mini-gin/flow/cmn"
	"github.com/korableg/mini-gin/flow/errs"
	"github.com/korableg/mini-gin/flow/msgs"
	"github.com/korableg/mini-gin/flow/node"
	"github.com/korableg/mini-gin/flow/repo"
)

type Hub struct {
	name  string
	nodes *node.Repository
}

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

func (h *Hub) Name() string {
	return h.name
}

func (h *Hub) AddNode(n *node.Node) error {
	return h.nodes.Store(n)
}

func (h *Hub) DeleteNode(n *node.Node) error {
	if n == nil {
		return nil
	}
	return h.nodes.Delete(n.Name())
}

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

func (h *Hub) deleteAllNodes() error {
	nodes := h.nodes.Slice()
	for _, n := range nodes {
		if err := h.DeleteNode(n); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) deleteNodeDB() error {
	err := h.nodes.DeleteDB()
	return err
}

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

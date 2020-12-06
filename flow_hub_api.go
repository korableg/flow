// Package flow implements basic methods to working with framework
package flow

import (
	"github.com/korableg/flow/errs"
	"github.com/korableg/flow/hub"
)

// NewHub creates new hub by name
func (m *Flow) NewHub(name string) (h *hub.Hub, err error) {

	if m.GetHub(name) != nil {
		err = errs.ErrHubIsAlreadyExists
		return
	}

	h, err = hub.New(name, m.db)
	err = m.hubs.Store(h)

	return
}

// GetHub gets the hub by name
func (m *Flow) GetHub(name string) (h *hub.Hub) {
	if value, ok := m.hubs.Load(name); ok {
		h = value
	}
	return
}

// GetAllHubs gets slice of hubs
func (m *Flow) GetAllHubs() []*hub.Hub {
	return m.hubs.Slice()
}

// AddNodeToHub adds node to hub by name
func (m *Flow) AddNodeToHub(hubName, nodeName string) error {

	h := m.GetHub(hubName)
	if h == nil {
		return errs.ErrHubNotFound
	}

	n := m.GetNode(nodeName)
	if n == nil {
		return errs.ErrNodeNotFound
	}

	return h.AddNode(n)

}

// DeleteNodeFromHub removes node from hub by name
func (m *Flow) DeleteNodeFromHub(hubName, nodeName string) error {

	h := m.GetHub(hubName)
	if h == nil {
		return errs.ErrHubNotFound
	}

	n := m.GetNode(nodeName)
	if n == nil {
		return errs.ErrNodeNotFound
	}

	return h.DeleteNode(n)

}

// DeleteHub removes hub
func (m *Flow) DeleteHub(name string) error {
	return m.hubs.Delete(name)
}

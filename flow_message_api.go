// Package flow implements basic methods to working with framework
package flow

import (
	"github.com/korableg/flow/errs"
	"github.com/korableg/flow/msgs"
)

// SendMessageToHub sends message from node to hub by name
func (m *Flow) SendMessageToHub(from, to string, data []byte) (*msgs.Message, error) {

	h := m.GetHub(to)
	if h == nil {
		return nil, errs.ErrHubNotFound
	}

	n := m.GetNode(from)
	if n == nil {
		return nil, errs.ErrNodeNotFound
	}

	mes := msgs.NewMessage(n.Name(), data)

	if err := h.PushMessage(mes); err != nil {
		return nil, err
	}

	return mes, nil

}

// SendMessageToNode sends message from node to node by name
func (m *Flow) SendMessageToNode(from, to string, data []byte) (*msgs.Message, error) {

	nodeTo := m.GetNode(to)
	if nodeTo == nil {
		return nil, errs.ErrNodeNotFound
	}

	nodeFrom := m.GetNode(from)
	if nodeFrom == nil {
		return nil, errs.ErrNodeNotFound
	}

	mes := msgs.NewMessage(nodeFrom.Name(), data)

	if err := nodeTo.PushMessage(mes); err != nil {
		return nil, err
	}

	return mes, nil

}

// GetMessage gets front message by node name
func (m *Flow) GetMessage(nodeName string) (*msgs.Message, error) {
	n := m.GetNode(nodeName)
	if n == nil {
		return nil, errs.ErrNodeNotFound
	}
	return n.FrontMessage(), nil
}

// RemoveMessage removes front message from node by name
func (m *Flow) RemoveMessage(nodeName string) error {
	n := m.GetNode(nodeName)
	if n == nil {
		return errs.ErrNodeNotFound
	}
	return n.RemoveFrontMessage()
}

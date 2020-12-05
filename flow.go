package flow

import (
	"github.com/korableg/mini-gin/flow/hub"
	"github.com/korableg/mini-gin/flow/node"
	"github.com/korableg/mini-gin/flow/repo"
)

type Flow struct {
	nodes *node.Repository
	hubs  *hub.Repository
	db    repo.DB
}

func New(db repo.DB) *Flow {

	var nodeDB repo.NodeDB
	if db != nil {
		nodeDB = db.NewNodeRepository()
	}

	f := new(Flow)
	f.db = db
	f.nodes = node.NewNodeRepository(nodeDB)
	f.hubs = hub.NewHubRepository(db, f.nodes)

	return f

}

func (m *Flow) Close() error {
	if err := m.nodes.Close(); err != nil {
		return err
	}
	if err := m.hubs.Close(); err != nil {
		return err
	}
	return nil
}

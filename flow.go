// Package flow implements basic methods to working with framework
package flow

import (
	"github.com/korableg/flow/hub"
	"github.com/korableg/flow/node"
	"github.com/korableg/flow/repo"
)

type Flow struct {
	nodes *node.Repository
	hubs  *hub.Repository
	db    repo.DB
}

// New creates new instance Flow. If db == nil all data keeps only memory and lost after restart
func New(db repo.DB) *Flow {

	var nodeDB repo.NodeDB
	if db != nil {
		nodeDB = db.NewNodeRepository()
	}

	f := new(Flow)
	f.db = db
	f.nodes = node.NewNodeRepository(nodeDB)
	f.hubs = hub.NewRepository(db, f.nodes)

	return f

}

// Close closes all dbs
func (m *Flow) Close() error {
	if err := m.nodes.Close(); err != nil {
		return err
	}
	if err := m.hubs.Close(); err != nil {
		return err
	}
	return nil
}

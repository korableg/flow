package hub

import (
	"github.com/korableg/mini-gin/flow/node"
	"github.com/korableg/mini-gin/flow/repo"
	"sync"
)

type Repository struct {
	hubs *sync.Map
	db   repo.HubDB
}

func NewHubRepository(db repo.DB, nodes *node.Repository) *Repository {

	hr := new(Repository)
	hr.hubs = new(sync.Map)

	if db == nil {
		return hr
	}

	hr.db = db.NewHubRepository()
	hubsRepo, err := hr.db.All()
	if err != nil {
		panic(err)
	}
	for _, hubRepo := range hubsRepo {
		hub, err := New(hubRepo.Name, db, nodes)
		if err != nil {
			panic(err)
		}
		hr.hubs.Store(hub.Name(), hub)
	}

	return hr
}

func (hr *Repository) Store(hub *Hub) error {
	if hr.db != nil {
		hubRepo := new(repo.Hub)
		hubRepo.Name = hub.Name()
		if err := hr.db.Store(hubRepo); err != nil {
			return err
		}
	}
	hr.hubs.Store(hub.Name(), hub)
	return nil
}

func (hr *Repository) Load(name string) (*Hub, bool) {
	if n, ok := hr.hubs.Load(name); ok {
		return n.(*Hub), ok
	}
	return nil, false
}

func (hr *Repository) Slice() []*Hub {
	hubs := make([]*Hub, 0, 20)
	f := func(value *Hub) error { hubs = append(hubs, value); return nil }
	_ = hr.Range(f)
	return hubs
}

func (hr *Repository) Range(f func(hub *Hub) error) error {
	var err error
	rangeFunc := func(key, value interface{}) bool {
		err = f(value.(*Hub))
		if err != nil {
			return false
		}
		return true
	}
	hr.hubs.Range(rangeFunc)
	return err
}

func (hr *Repository) Delete(name string) (err error) {

	if hub, ok := hr.Load(name); ok {
		if err = hub.deleteNodeDB(); err != nil {
			return
		}
		if hr.db != nil {
			if err = hr.db.Delete(name); err != nil {
				return err
			}
		}
		hr.hubs.Delete(name)
	}
	return

}

func (hr *Repository) Close() error {
	if hr.db == nil {
		return nil
	}
	return hr.db.Close()
}

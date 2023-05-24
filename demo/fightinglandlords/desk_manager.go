package fightinglandlords

import "sync"

// 桌子管理器
type deskManager struct {
	lock  sync.Mutex
	dmMap map[string]*Desk
}

var dkMgr *deskManager

func init() {
	dkMgr = &deskManager{
		dmMap: make(map[string]*Desk),
	}
}

func (m *deskManager) Add(d *Desk) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.dmMap[d.Id] = d
}

func (m *deskManager) Del(d *Desk) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.dmMap, d.Id)
}

func (m *deskManager) Get(deskId string) (*Desk, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	desk, ok := m.dmMap[deskId]
	if !ok {
		return nil, DeskNotFound
	}
	return desk, nil
}

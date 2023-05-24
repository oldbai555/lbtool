package fightinglandlords

import "sync"

// 桌子管理器
type playerManager struct {
	lock  sync.Mutex
	pyMap map[string]*Player
}

var pyMgr *playerManager

func init() {
	pyMgr = &playerManager{
		pyMap: make(map[string]*Player),
	}
}

func (m *playerManager) Add(d *Player) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.pyMap[d.Id] = d
}

func (m *playerManager) Del(d *Player) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.pyMap, d.Id)
}

func (m *playerManager) Get(pyId string) (*Player, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	desk, ok := m.pyMap[pyId]
	if !ok {
		return nil, DeskNotFound
	}
	return desk, nil
}

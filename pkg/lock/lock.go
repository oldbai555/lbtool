package lock

import "sync"

type St struct {
	sync.RWMutex
}

func NewWithUsageMu(internalMu sync.Locker, onUsageAdd func(), onUsageSub func()) *WithUsageMu {
	return &WithUsageMu{
		onUsageAdd: onUsageAdd,
		onUsageSub: onUsageSub,
		internalMu: internalMu,
	}
}

type WithUsageMu struct {
	outerMu    sync.RWMutex
	internalMu sync.Locker
	UsageNum   uint32
	onUsageSub func()
	onUsageAdd func()
}

func (m *WithUsageMu) Lock() {
	m.internalMu.Lock()
	m.UsageNum++
	m.onUsageAdd()
	m.internalMu.Unlock()

	m.outerMu.Lock()
}

func (m *WithUsageMu) Unlock() {
	m.internalMu.Lock()
	m.UsageNum--
	m.onUsageSub()
	m.internalMu.Unlock()

	m.outerMu.Unlock()
}

func (m *WithUsageMu) RLock() {
	m.internalMu.Lock()
	m.UsageNum++
	m.onUsageAdd()
	m.internalMu.Unlock()

	m.outerMu.RLock()
}

func (m *WithUsageMu) RUnlock() {
	m.internalMu.Lock()
	m.UsageNum--
	m.onUsageSub()
	m.internalMu.Unlock()

	m.outerMu.RUnlock()
}

type MulElemMuFactory struct {
	elemMuMap      map[interface{}]*WithUsageMu
	opMapMu        sync.Mutex
	makeOrGetMapMu sync.Mutex
}

func NewMulElemMuFactory() *MulElemMuFactory {
	return &MulElemMuFactory{
		elemMuMap: map[interface{}]*WithUsageMu{},
	}
}

func (m *MulElemMuFactory) MakeOrGetSpecElemMu(elem interface{}) *WithUsageMu {
	m.makeOrGetMapMu.Lock()
	defer m.makeOrGetMapMu.Unlock()
	mu, ok := m.elemMuMap[elem]
	if !ok {
		mu = NewWithUsageMu(
			&m.opMapMu,
			func() {
				// save this lock while any thread used the handler lock
				if _, ok := m.elemMuMap[elem]; ok {
					return
				}
				m.elemMuMap[elem] = mu
			},
			func() {
				// remove this lock from map if this lock might have no owner
				if mu.UsageNum > 0 {
					return
				}
				delete(m.elemMuMap, elem)
			},
		)

		m.elemMuMap[elem] = mu
	}
	return mu
}

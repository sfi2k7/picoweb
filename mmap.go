package picoweb

import "sync"

//Hub Code
type mmap struct {
	m     map[string]*handler
	_lock *sync.Mutex
}

func (m *mmap) add(id string, h *handler) {
	m._lock.Lock()
	defer m._lock.Unlock()

	m.m[id] = h
}

func (m *mmap) remove(id string) error {
	m._lock.Lock()
	defer m._lock.Unlock()

	delete(m.m, id)

	return nil
}

func (m *mmap) count() int {
	m._lock.Lock()
	defer m._lock.Unlock()

	return len(m.m)
}

func (m *mmap) get(id string) *handler {
	m._lock.Lock()
	defer m._lock.Unlock()

	h, ok := m.m[id]
	if !ok {
		return nil
	}
	return h
}

func (m *mmap) closeAll() {
	m._lock.Lock()
	defer m._lock.Unlock()

	if m.count() == 0 {
		return
	}

	for _, h := range m.m {
		if h == nil || h.isOpen == false {
			continue
		}

		h.forceExit()
	}
}

func newmmap() *mmap {
	return &mmap{
		m:     make(map[string]*handler),
		_lock: &sync.Mutex{},
	}
}

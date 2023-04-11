package picoweb

import "sync"

//Hub Code
type genericmmap struct {
	m     map[string]*genericconnectionhandler
	_lock *sync.Mutex
}

//add adds entry.
//Add a new entry to the map
func (m *genericmmap) add(id string, h *genericconnectionhandler) {
	m._lock.Lock()
	defer m._lock.Unlock()

	m.m[id] = h
}

func (m *genericmmap) remove(id string) error {
	m._lock.Lock()
	defer m._lock.Unlock()

	delete(m.m, id)

	return nil
}

func (m *genericmmap) count() int {
	m._lock.Lock()
	defer m._lock.Unlock()

	return len(m.m)
}

func (m *genericmmap) get(id string) *genericconnectionhandler {
	m._lock.Lock()
	defer m._lock.Unlock()

	h, ok := m.m[id]
	if !ok {
		return nil
	}
	return h
}

func (m *genericmmap) closeAll() {
	m._lock.Lock()
	defer m._lock.Unlock()

	if m.count() == 0 {
		return
	}

	for _, h := range m.m {
		if h == nil || !h.isOpen {
			continue
		}

		h.Terminate()
	}
}

func (m *genericmmap) send(id string, data WsData) {
	m._lock.Lock()
	defer m._lock.Unlock()

	h, ok := m.m[id]
	if !ok || !h.isOpen {
		return
	}

	h.out.In(data)
}

func (m *genericmmap) broadcast(data WsData) {
	m._lock.Lock()
	defer m._lock.Unlock()

	for _, h := range m.m {
		if h == nil || !h.isOpen {
			continue
		}

		h.out.In(data)
	}
}

func newgenericmmap() *genericmmap {
	return &genericmmap{
		m:     make(map[string]*genericconnectionhandler),
		_lock: &sync.Mutex{},
	}
}

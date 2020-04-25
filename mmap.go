package picoweb

import "sync"

//Hub Code
type mmap struct {
	m map[string]*handler

	c2u   map[string]string   // Connection -> User
	u2cs  map[string][]string // User -> []Connections
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
	//m.logoff(id)

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

func (m *mmap) logoff(cid string) {

	u, ok := m.c2u[cid]
	if !ok {
		return
	}

	cs, ok := m.u2cs[u]
	if !ok {
		return // Not Logged in
	}

	var newList []string
	for _, c := range cs {
		if c == cid {
			continue
		}
		newList = append(newList, c)
	}

	if len(newList) == 0 { // No Connection
		delete(m.u2cs, u) //Delete Orphaned Entry
		return
	}

	m.u2cs[u] = newList
}

func (m *mmap) login(u, cid string) {
	m._lock.Lock()
	defer m._lock.Unlock()

	m.c2u[cid] = u
	_, ok := m.u2cs[u]
	if !ok {
		m.u2cs[u] = []string{}
	}
	m.u2cs[u] = append(m.u2cs[u], cid)
}

func newmmap() *mmap {
	return &mmap{
		m:     make(map[string]*handler),
		c2u:   make(map[string]string),
		u2cs:  make(map[string][]string),
		_lock: &sync.Mutex{},
	}
}

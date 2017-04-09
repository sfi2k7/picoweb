package picoweb

import (
	"sync"
)

var (
	flashMutex sync.Mutex
)

func init() {
	flashMutex = sync.Mutex{}
}

type Flash map[string]interface{}

func (f Flash) Get(sessionId string) interface{} {
	flashMutex.Lock()
	defer flashMutex.Unlock()

	v, ok := f[sessionId]
	if ok {
		delete(f, sessionId)
	}

	return v
}

func (f Flash) Set(sessionId string, value interface{}) {
	flashMutex.Lock()
	defer flashMutex.Unlock()

	f[sessionId] = value
}

func (f Flash) Has(sessionId string) bool {
	_, ok := f[sessionId]
	return ok
}

func (f Flash) Clear() {
	flashMutex.Lock()
	defer flashMutex.Unlock()

	for k, _ := range f {
		delete(f, k)
	}
}

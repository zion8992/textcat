package tc

import (
	"sync"
)

type SessionManager struct {
	Mu    sync.RWMutex
	Store map[int]any
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Store: make(map[int]any),
	}
}

func (sm *SessionManager) Add(id int, object any) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()
	sm.Store[id] = object
 
}

func (sm *SessionManager) Remove(id int) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	delete(sm.Store, id)

}

func (sm *SessionManager) GetUnused() int {
	sm.Mu.RLock()
	defer sm.Mu.RUnlock()

	max := 0
	for id := range sm.Store {
		if id > max {
			max = id
		}
	}

	return max + 1
}


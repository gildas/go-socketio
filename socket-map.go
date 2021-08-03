package socketio

import "sync"

type socketMap struct {
	sockets map[string]Socket
	lock    sync.RWMutex
}

func NewSocketMap() *socketMap {
	return &socketMap{
		sockets: map[string]Socket{},
		lock:    sync.RWMutex{},
	}
}

func (sm *socketMap) Add(socket Socket) *socketMap {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.sockets[socket.ID()] = socket
	return sm
}

func (sm *socketMap) Get(id string) Socket {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return sm.sockets[id]
}

func (sm *socketMap) Delete(id string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.sockets, id)
}

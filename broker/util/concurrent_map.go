package util

import (
	"broker/proto"
	"sync"
)

type ConcurrentMapData struct {
	Mux sync.RWMutex
	Map map[string]chan *proto.Data
}

func (t *ConcurrentMapData) Read(key string) (chan *proto.Data, bool) {
	t.Mux.RLock()
	defer t.Mux.RUnlock()
	val, ok := t.Map[key]
	return val, ok
}

func (t *ConcurrentMapData) Write(key string, value chan *proto.Data) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	t.Map[key] = value
}

func (t *ConcurrentMapData) Delete(key string) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	delete(t.Map, key)
}



type ConcurrentMapBool struct {
	Mux sync.RWMutex
	Map map[string]chan bool
}

func (t *ConcurrentMapBool) Read(key string) (chan bool, bool) {
	t.Mux.RLock()
	defer t.Mux.RUnlock()
	val, ok := t.Map[key]
	return val, ok
}

func (t *ConcurrentMapBool) Write(key string, value chan bool) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	t.Map[key] = value
}

func (t *ConcurrentMapBool) Delete(key string) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	delete(t.Map, key)
}
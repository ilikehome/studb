package index

import (
	"sync"
	"encoding/hex"
)

type memIndex struct{
	inx map[string]int64
	lock *sync.RWMutex
}

func createMemIndex(m map[string]int64) *memIndex {
	mi := new(memIndex)
	mi.inx = m
	mi.lock = new(sync.RWMutex)
	return mi
}

func (m *memIndex) put(k []byte, locate int64) {
	m.lock.Lock()
	defer  m.lock.Unlock()

	m.inx[hex.EncodeToString(k)] = locate
}

func (m *memIndex) get(k []byte) (int64, bool){
	m.lock.RLock()
	defer  m.lock.RUnlock()

	v,ok := m.inx[hex.EncodeToString(k)]
	return v,ok
}

func (m *memIndex) del(k []byte){
	m.lock.Lock()
	defer  m.lock.Unlock()

	delete(m.inx, hex.EncodeToString(k))
}

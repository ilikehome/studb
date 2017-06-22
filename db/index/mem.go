package index

import (
	"sync"
	"encoding/hex"
	"os"
)

type indexInMem struct{
	inx map[string]int64
	lock sync.RWMutex
}

func createMemIndex(m map[string]int64) *indexInMem {
	mi := new(indexInMem)
	mi.inx = m
	return mi
}

func load(f *os.File) map[string]int64{//TODO:read data file when start, now.
	inx := make(map[string]int64)
	buf := [290]byte{}
	i := int64(0)
	for {
		_, err := f.ReadAt(buf[:], i)
		if err != nil{
			break
		}
		inx[hex.EncodeToString(buf[2 : 2+buf[0]])] = i
		i += 290
	}
	return inx
}

func (m *indexInMem) put(k []byte, locate int64) {
	m.inx[hex.EncodeToString(k)] = locate
}

func (m *indexInMem) get(k []byte) (int64, bool){
	v,ok := m.inx[hex.EncodeToString(k)]
	return v,ok
}

func (m *indexInMem) del(k []byte){
	delete(m.inx, hex.EncodeToString(k))
}

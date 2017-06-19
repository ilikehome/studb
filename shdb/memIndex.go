package shdb

import (
	"sync"
	"encoding/hex"
	"os"
)


type memInx struct{
	inx map[string]int64
	lock sync.RWMutex
	dumpInDisk *os.File
}

func createMemInx(f *os.File) *memInx{
	mi := new(memInx)
 	inx := load(f)
	mi.inx = inx
	return mi
}

func load(f *os.File) map[string]int64{
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

func (m *memInx) flush() {

}

func (m *memInx)put(k []byte, locate int64) {
	m.inx[hex.EncodeToString(k)] = locate
}

func (m *memInx)get(k []byte) (int64, bool){
	v,ok := m.inx[hex.EncodeToString(k)]
	return v,ok
}

func (m *memInx)del(k []byte){
	delete(m.inx, hex.EncodeToString(k))
}

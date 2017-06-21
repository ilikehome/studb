package index

import (
	"sync"
	"encoding/hex"
	"os"
)

type IndexInMem struct{
	inx map[string]int64
	lock sync.RWMutex
}

func CreateMemInx(f *os.File) *IndexInMem{
	mi := new(IndexInMem)
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

func (m *IndexInMem) Flush() {

}

func (m *IndexInMem)Put(k []byte, locate int64) {
	m.inx[hex.EncodeToString(k)] = locate
}

func (m *IndexInMem)Get(k []byte) (int64, bool){
	v,ok := m.inx[hex.EncodeToString(k)]
	return v,ok
}

func (m *IndexInMem)Del(k []byte){
	delete(m.inx, hex.EncodeToString(k))
}
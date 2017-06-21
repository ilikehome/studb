package index

import "os"

type Index struct{
	iid *indexInDisk
	m *indexInMem
}

func Init(diskFile *os.File) *Index{
	i := new(Index)
	i.m = createMemIndex(diskFile)
	return i
}

func (i *Index)Put(k []byte, seq uint64, locate int64){
	i.m.put(k, locate)
}

func (i *Index)Get(k []byte) (int64, bool){
	return i.m.get(k)
}

func (i *Index)Del(k []byte){
	i.m.del(k)
}

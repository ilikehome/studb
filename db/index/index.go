package index

import (
	"os"
	"github.com/ilikehome/studb/db/constant"
)

type Index struct{
	iid *indexInDisk
	m *indexInMem
}

func Init(f *os.File) *Index{
	i := new(Index)
	i.iid = open(f)
	i.m = createMemIndex(i.iid.readToMem())
	return i
}

func (i *Index)Put(seq uint64, k []byte, loc int64){
	i.m.put(k, loc)
	i.iid.append(seq, constant.OP_PUT , k, uint64(loc))
}

func (i *Index)Get(k []byte) (int64, bool){
	return i.m.get(k)
}

func (i *Index)Del(seq uint64, k []byte){
	i.m.del(k)
	i.iid.append(seq, constant.OP_DEL , nil, 0)
}

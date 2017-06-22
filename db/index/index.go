package index

import (
	"os"
	"github.com/ilikehome/studb/db/constant"
)

type Index struct{
	di *diskIndex
	mi *memIndex
}

func Init(f *os.File) *Index{
	i := new(Index)
	i.di = openDiskIndex(f)
	i.mi = createMemIndex(i.di.readToMem())
	return i
}

func (i *Index)Put(seq uint64, k []byte, loc int64){
	i.mi.put(k, loc)
	i.di.append(seq, constant.OP_PUT , k, uint64(loc))
}

func (i *Index)Get(k []byte) (int64, bool){
	return i.mi.get(k)
}

func (i *Index)Del(seq uint64, k []byte){
	i.mi.del(k)
	i.di.append(seq, constant.OP_DEL , nil, 0)
}

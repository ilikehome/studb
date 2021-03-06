package index

import(
	"encoding/binary"
	"os"
	"fmt"
	"sync"
	"io"
	"encoding/hex"
	"github.com/ilikehome/studb/db/constant"
)

type diskIndex struct{
	lock *sync.Mutex
	f *os.File
}

type record struct{
	seq uint64//8
	op constant.OPT_CRUD//1
	klen uint16//2
	k []byte//kLen
	loc uint64//8
}

func openDiskIndex(f *os.File) *diskIndex {
	di := new(diskIndex)
	di.f = f
	di.lock = new(sync.Mutex)
	return di
}

func (di *diskIndex) readToMem() map[string]int64{
	di.lock.Lock()
	defer di.lock.Unlock()

	di.f.Seek(0, os.SEEK_SET)
	defer di.f.Seek(0, os.SEEK_END)

	mem := make(map[string]int64)
	buf := [1024*1024*64]byte{}
	for{
		var leftbuf []byte
		n, err := di.f.Read(buf[:])
		if err==io.EOF{
			return mem
		}else if err != nil{
			panic(fmt.Sprintf("open index file to mem fail. %v", err))
		}else{
			if leftbuf == nil{
				offset := 0
				for offset+1 < n{
					kSize := int(binary.BigEndian.Uint16(buf[offset+9 : offset +11]))
					kStr := hex.EncodeToString(buf[offset+11 : offset+11+kSize])
					loc := int64(binary.BigEndian.Uint64(buf[offset+11+kSize : offset+19+kSize]))
					mem[kStr] = loc
					offset = offset+19+kSize
				}
				return mem
			}//TODO: more than 64M is not support
		}
	}

}

func (di *diskIndex)compact() error{
	di.lock.Lock()
	defer  di.lock.Unlock()

	return nil
}

func (i *record)toBytes() []byte{
	size := 8+1+2+i.klen+8
	b := make([]byte, size, size)
	binary.BigEndian.PutUint64(b[0:8], i.seq)
	b[8]=byte(i.op)
	binary.BigEndian.PutUint16(b[9:11], i.klen)
	copy(b[11:11+i.klen], i.k[:])
	binary.BigEndian.PutUint64(b[11+i.klen:], i.loc)
	return b
}

func (di *diskIndex) append(seq uint64, op constant.OPT_CRUD,  k []byte, loc uint64) error{
	di.lock.Lock()
	defer di.lock.Unlock()

	i := record{seq, op, uint16(len(k)), k, loc}
	di.f.Write(i.toBytes())
	return nil
}

func (di *diskIndex) close(){
	di.f.Close()
}



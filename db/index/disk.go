package index

import(
	"encoding/binary"
	"os"
	"fmt"
	"sync"
	"io"
	"encoding/hex"
)

type indexInDisk struct{
	lock sync.Mutex
	f *os.File
}

type record struct{
	kLen uint16
	loc uint64
	seq uint64
	k []byte
}

func Open(diskFile string) *indexInDisk{
	f,err := os.OpenFile(diskFile, os.O_RDWR, 0666);
	if err != nil{
		panic(fmt.Sprintf("open index file in disk fail. %v", err))
	}
	iid := new(indexInDisk)
	iid.f = f
	return iid
}

func (iid *indexInDisk) readToMem() map[string]int64{
	iid.lock.Lock()
	defer iid.lock.Unlock()

	iid.f.Seek(0, os.SEEK_SET)
	defer iid.f.Seek(0, os.SEEK_END)

	mem := make(map[string]int64)
	buf := [1024*1024*64]byte{}
	for{
		var leftbuf []byte
		n, err := iid.f.Read(buf[:])
		if err==io.EOF{
			return mem
		}else if err != nil{
			panic(fmt.Sprintf("open index file to mem fail. %v", err))
		}else{
			if leftbuf == nil{
				offset := 0
				for offset+1 < n{
					kSize := int(binary.BigEndian.Uint16(buf[offset:offset +2]))
					loc := int64(binary.BigEndian.Uint64(buf[offset+2:offset+10]))
					kStr := hex.EncodeToString(buf[offset+18:offset+18+kSize])
					mem[kStr] = loc
					offset = offset+18+kSize
				}
				return mem
			}//TODO:
		}
	}

}

func (iid *indexInDisk)compact() error{
	return nil
}

func (i *record)toBytes() []byte{
	size := 2+8+8+i.kLen
	b := make([]byte, size, size)
	binary.BigEndian.PutUint16(b[0:2], i.kLen)
	binary.BigEndian.PutUint64(b[2:10], i.loc)
	binary.BigEndian.PutUint64(b[10:18], i.seq)
	copy(b[18:18+i.kLen], i.k[:])
	return b
}

func (iid *indexInDisk) append(loc uint64, seq uint64, k []byte) error{
	iid.lock.Lock()
	defer iid.lock.Unlock()

	i := record{uint16(len(k)), loc, seq, k}
	iid.f.Write(i.toBytes())
	return nil
}

func (iid *indexInDisk) close(){
	iid.f.Close()
}



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
	fname string
}

type record struct{
	klen uint16
	loc int64
	k []byte

}

func Open(diskFile string) *indexInDisk{
	f,err := os.OpenFile(diskFile, os.O_APPEND, 0666);
	if err != nil{
		panic(fmt.Sprintf("open index file in disk fail. %v", err))
	}
	iid := new(indexInDisk)
	iid.f = f
	iid.fname = diskFile
	return iid
}

func (iid *indexInDisk)ReadToMem() map[string]int64{
	iid.lock.Lock()
	defer iid.lock.Unlock()

	mem := make(map[string]int64)
	if f,err := os.OpenFile(iid.fname, os.O_RDONLY, 0666); err!=nil{
		buf := [1024*1024*64]byte{}
		var leftbuf []byte
		for{//TODO:more than 64M
			n, err := f.Read(buf[:])
			if err==io.EOF{
				return mem
			}else if err != nil{
				panic(fmt.Sprintf("open index file to mem fail. %v", err))
			}else{
				if leftbuf == nil{
					offset := 0
					for offset+1<n {
						kLen := int(binary.BigEndian.Uint16(buf[offset:offset +2]))
						loc := int64(binary.BigEndian.Uint64(buf[offset+2:offset+10]))
						keyStr := hex.EncodeToString(buf[offset+10:offset+10+kLen])
						mem[keyStr] = loc
						offset = offset+10+kLen
					}
				}//TODO:
				return mem
			}
		}

	}else{
		panic(fmt.Sprintf("open index file in disk fail. %v", err))
	}

}

func (iid *indexInDisk)compact() error{
	return nil
}

func (i *record)toBytes() []byte{
	size := 2+8+i.klen
	b := make([]byte, size, size)
	binary.BigEndian.PutUint16(b[0:2], i.klen)
	binary.BigEndian.PutUint64(b[2:10], uint64(i.loc))
	copy(b[10:], i.k[:])
	return b
}

func (iid *indexInDisk)Append(k []byte, loc int64) error{
	i := record{uint16(len(k)), loc, k}
	iid.f.Write(i.toBytes())
	return nil
}



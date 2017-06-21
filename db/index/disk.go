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
	kLen uint16
	loc uint64
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
		for{//TODO: larger than 64M
			n, err := f.Read(buf[:])
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
						kStr := hex.EncodeToString(buf[offset+12:offset+12+kSize])
						mem[kStr] = loc
						offset = offset+12+kSize
					}
				return mem
				}//TODO:
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
	size := 2+8+i.kLen
	b := make([]byte, size, size)
	binary.BigEndian.PutUint16(b[0:2], i.kLen)
	binary.BigEndian.PutUint64(b[4:12], i.loc)
	copy(b[12:12+i.kLen], i.k[:])
	return b
}

func (iid *indexInDisk)Append(k []byte, loc uint64) error{
	i := record{uint16(len(k)), loc, k}
	iid.f.Write(i.toBytes())
	return nil
}



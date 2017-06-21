package index

import(
	"encoding/binary"
	"os"
	"fmt"
	"sync"
	"io"
)

type indexInDisk struct{
	lock sync.Mutex
	f *os.File
	fname string
}

type record struct{
	k string
	location int64
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
		for{
			_, err := f.Read(buf[:])
			if err==io.EOF{
				return mem
			}else if err != nil{
				panic(fmt.Sprintf("open index file to mem fail. %v", err))
			}else{
				if leftbuf == nil{
					offset := 0
					for {
						kSize := int(binary.BigEndian.Uint16(buf[offset:offset +2]))
						seq := binary.BigEndian.Uint64(buf[offset+4:offset+12])
						k := buf[offset+12:offset+12+kSize]
						v := buf[offset+12+kSize:offset+12+kSize+vSize]
						mem.put()
					}

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
	size := 2+2+8+i.kSize+i.vSize
	b := make([]byte, size, size)
	binary.BigEndian.PutUint16(b[0:2], i.kSize)
	binary.BigEndian.PutUint16(b[2:4], i.vSize)
	binary.BigEndian.PutUint64(b[4:12], i.seq)
	copy(b[12:12+i.kSize], i.k[:])
	binary.BigEndian.PutUint64(b[4:12], i.seq)
	copy(b[12+i.kSize:], i.v[:])
	return b
}

func (iid *indexInDisk)Append(seq uint64, k,v []byte) error{
	i := record{uint16(len(k)), uint16(len(v)), seq, k, v}
	iid.f.Write(i.toBytes())
	return nil
}



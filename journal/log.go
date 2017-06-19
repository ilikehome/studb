package journal

import (
	"os"
	"sync"
	"github.com/ilikehome/studb/shdb"
)

const(
	blockSize = 32*1024
)

type RECORD_POSITION int8

const(
	fullChunkType RECORD_POSITION = iota
	firstChunkType
	middleChunkType
	lastChunkType
)

const ChunkMaxSize = 32*1024

type Log struct{
	f *os.File
	lock sync.Mutex
	buf [ChunkMaxSize]byte
	bufCnt int64
}

type record struct{
	len int64
	seq int64
	pos RECORD_POSITION
	op shdb.OP_TYPE
	content []byte
}

func (r *record)getBytes() []byte{

}

func OpenJournal(journal string) *Log {
	l := new(Log)
	jf,_ := os.OpenFile(journal, os.O_APPEND, 0666)
	l.f = jf
	if fInfo,err := jf.Stat(); err==nil{
		l.bufCnt = fInfo.Size() % 32*1024
	}
	l.buf = [32*1024]byte{}
	return l
}

func (l *Log)Write(batch *[]shdb.Row) error{
	for _,r := range *batch{
		size := r.KLen+r.VLen + 8 + 1 + 8 +1
		if (l.bufCnt + int64(size)) <= ChunkMaxSize{
			record := new(record)
			record.len = int64(size)
			record.seq = r.Seq
			record.pos = fullChunkType
			record.op = shdb.OP_PUT
			c := [r.KLen+r.VLen]byte{}
			copy(c[:r.KLen], r.KeyValue[:r.KLen])
			copy(c[r.KLen:], r.KeyValue[32:32+r.VLen])
			record.content = c[:]
			l.f.Write(record.getBytes())
		}else{

		}
	}
	return nil
}



func (l *Log)Close() {
	l.f.Close()
}

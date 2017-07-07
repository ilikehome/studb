package wal

import (
	"os"
	"sync"
	"encoding/binary"
	"github.com/ilikehome/studb/db/constant"
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

type row struct{
	Seq int64
	KLen, VLen uint8
	KeyValue [290]byte//1+1+32+256
}

type Log struct{
	f           *os.File
	lock        sync.Mutex
	chunkOffset int64
}

type record struct{
	len     int64
	seq     int64
	pos     RECORD_POSITION
	op      constant.OPT_CRUD
	content []byte
}

func OpenJournal(journal string) *Log {
	l := new(Log)
	jf,_ := os.OpenFile(journal, os.O_APPEND, 0666)
	l.f = jf
	if fInfo,err := jf.Stat(); err==nil{
		l.chunkOffset = fInfo.Size() % 32*1024
	}
	return l
}

func (l *Log)Write(batch *[]row) error{
	for _,r := range *batch{
		size := r.KLen+r.VLen + 8 + 1 + 8 +1
		if (l.chunkOffset + int64(size)) <= ChunkMaxSize{
			record := new(record)
			record.len = int64(size)
			record.seq = r.Seq
			record.pos = fullChunkType
			record.op = constant.OP_PUT
			c := make([]byte, size, size)
			var buf = make([]byte, 8)
			binary.BigEndian.PutUint64(buf, uint64(record.len))
			copy(c[:8], buf)
			binary.BigEndian.PutUint64(buf, uint64(record.seq))
			copy(c[8:16], buf)
			c[16]= byte(record.pos)
			c[17]= byte(record.op)
			copy(c[18:18+r.KLen], r.KeyValue[:r.KLen])
			copy(c[18+r.KLen:], r.KeyValue[32:32+r.VLen])
			record.content = c[:]
			l.f.Write(c[:])
			l.chunkOffset+= int64(size)
		}else{

		}
	}
	return nil
}

func (l *Log)Close() {
	l.f.Close()
}
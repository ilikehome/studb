package journal

import (
	"os"
	"sync"
	"github.com/ilikehome/studb/shdb"
)

const(
	blockSize = 32*1024
)

type RECORD_POSITION int

const(
	fullChunkType RECORD_POSITION = iota
	firstChunkType
	middleChunkType
	lastChunkType
)

type log struct{
	f *os.File
	lock sync.Mutex
	buf [32*1024]byte
	bufCnt int
}

type record struct{
	seq int64
	len int64
	pos RECORD_POSITION
	content []byte
}

func OpenJournal() *os.File{
	f,_ := os.OpenFile(dbFile, os.O_APPEND, 0666)
	return f
}

func (l *log)Write(batch *[]shdb.Row) error{

	return nil
}

func (l *log)Close() {
	l.f.Close()
}

func main() {
}

package studb

import (
	"sync"
	"os"
	"github.com/ilikehome/studb/journal"
	"github.com/ilikehome/studb/index"
)

type DB struct{
	lock sync.RWMutex
	seq int64
	diskFile *os.File
	mi *index.IndexInMem
	j *journal.Log
}

func Load(dbFile string ) *DB{
	f,_ := os.OpenFile(dbFile, os.O_RDWR, 0666)
	db := new(DB)
	mi := index.CreateMemInx(f)
	db.diskFile = f
	db.mi = mi
	db.j = journal.OpenJournal(dbFile+".j")
	return db
}

type Row struct{
	Seq int64
	KLen, VLen uint8
	KeyValue [290]byte//1+1+32+256
}

func (db *DB) write(r *Row, locate int64) error{
	db.diskFile.Seek(locate, os.SEEK_SET)
	_, err := db.diskFile.Write(append([]byte{}, r.KeyValue[:]...))
	return err
}

func (db *DB) writeEnd(r *Row) error{
	db.diskFile.Seek(0, os.SEEK_END)
	_, err := db.diskFile.Write(append([]byte{}, r.KeyValue[:]...))
	return err
}

func (db *DB) Write(k,v []byte) error{
	inx,ok := db.mi.Get(k)
	r := new(Row)
	r.KLen = uint8(len(k))
	r.VLen = uint8(len(v))
	r.KeyValue[0] = r.KLen
	r.KeyValue[1] = r.VLen
	copy(r.KeyValue[2:33], k)
	copy(r.KeyValue[33:], v)
	if ok{
		db.mi.Put(k, inx)
		return db.write(r, int64(inx))
	}else{
		fi,_ := db.diskFile.Stat()
		db.mi.Put(k, fi.Size())
		return db.writeEnd(r)
	}
}

func (db *DB) Read(k []byte) []byte{
	inx,ok := db.mi.Get(k)
	if !ok{
		return nil
	}
	buf := [290]byte{}
	db.diskFile.ReadAt(buf[:], int64(inx))
	return buf[33 : 33+buf[1]]
}

func (db *DB) Close(){
	db.diskFile.Close()
}



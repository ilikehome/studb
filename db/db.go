package db

import (
	"sync"
	"os"
	"github.com/ilikehome/studb/db/journal"
	"github.com/ilikehome/studb/db/index"
	"fmt"
)

type DB struct{
	lock sync.RWMutex
	seq int64
	diskFile *os.File
	mi *index.IndexInMem
	j *journal.Log
}

func Open(dbFile string ) *DB{
	f,_ := os.OpenFile(dbFile, os.O_RDWR, 0666)
	db := new(DB)
	mi := index.CreateMemInx(f)
	db.diskFile = f
	db.mi = mi
	db.j = journal.OpenJournal(dbFile+".j")
	return db
}

func (db *DB) write(r *[290]byte, locate int64) error{
	db.diskFile.Seek(locate, os.SEEK_SET)
	_, err := db.diskFile.Write(r[:])
	return err
}

func (db *DB) writeEnd(r *[290]byte) error{
	db.diskFile.Seek(0, os.SEEK_END)
	_, err := db.diskFile.Write(r[:])
	return err
}

func (db *DB) Write(k,v []byte) error{
	inx,ok := db.mi.Get(k)
	kv := [290]byte{}//1+1+32+256
	kv[0] = uint8(len(k))
	kv[1] = uint8(len(v))
	copy(kv[2:33], k)
	copy(kv[33:], v)
	if ok{
		db.mi.Put(k, inx)
		return db.write(&kv, int64(inx))
	}else{
		fi,_ := db.diskFile.Stat()
		db.mi.Put(k, fi.Size())
		return db.writeEnd(&kv)
	}
}

func (db *DB) Read(k []byte) ([]byte, error){
	inx,ok := db.mi.Get(k)
	if !ok{
		return nil, fmt.Errorf("DB is not ok.")
	}
	buf := [290]byte{}
	db.diskFile.ReadAt(buf[:], int64(inx))
	return buf[33 : 33+buf[1]], nil
}

func (db *DB) Close(){
	db.diskFile.Close()
}


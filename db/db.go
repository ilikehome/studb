package db

import (
	"sync"
	"os"
	"github.com/ilikehome/studb/db/wal"
	"github.com/ilikehome/studb/db/index"
	"github.com/ilikehome/studb/db/cache"
	"fmt"
)

const(
	DB_SUFFIX = ".db"
	INX_SUFFIX = ".inx"
)
type DB struct{
	lock     sync.RWMutex
	seq      int64
	diskFile *os.File
	inx      *index.Index
	j        *wal.Log
	c		*cache.Cache
}

func Open(dbName string ) *DB{
	dbFilePath := dbName + DB_SUFFIX
	inxFilePath := dbName + INX_SUFFIX
	f,err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil{
		panic(err.Error())
	}
	db := new(DB)

	f,err = os.OpenFile(inxFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil{
		panic(err.Error())
	}
	inx := index.Init(f)
	db.diskFile = f
	db.inx = inx
	db.j = wal.OpenJournal(dbName +".j")
	db.c = cache.CreateCache(10)
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
	//db.j.Write()
	db.c.Put(k,v)

	inx,ok := db.inx.Get(k)
	kv := [290]byte{}//1+1+32+256
	kv[0] = uint8(len(k))
	kv[1] = uint8(len(v))
	copy(kv[2:33], k)
	copy(kv[33:], v)
	if ok{
		db.inx.Put(1, k, inx)
		return db.write(&kv, int64(inx))
	}else{
		fi,_ := db.diskFile.Stat()
		db.inx.Put(1, k, fi.Size())
		return db.writeEnd(&kv)
	}
}

func (db *DB) Read(k []byte) ([]byte, error){
	v := db.c.Get(k)
	if v != nil{
		return v.([]byte), nil
	}

	inx,ok := db.inx.Get(k)
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



package shdb

import (
	"sync"
	"os"
)

type DB struct{
	lock sync.RWMutex
	diskFile *os.File
	mi *memInx
}
func Load(dbFile string ) *DB{
	f,_ := os.OpenFile(dbFile, os.O_RDWR, 0666)
	db := new(DB)
	mi := createMemInx(f)
	db.diskFile = f
	db.mi = mi
	return db
}

type row struct{
	kLen,vLen uint8
	KeyValue [290]byte//1+1+32+256
}

func (db *DB) write(r *row, locate int64) error{
	db.diskFile.Seek(locate, os.SEEK_SET)
	_, err := db.diskFile.Write(append([]byte{}, r.KeyValue[:]...))
	return err
}

func (db *DB) writeEnd(r *row) error{
	db.diskFile.Seek(0, os.SEEK_END)
	_, err := db.diskFile.Write(append([]byte{}, r.KeyValue[:]...))
	return err
}

func (db *DB) Write(k,v []byte) error{
	inx,ok := db.mi.get(k)
	r := new(row)
	r.kLen = uint8(len(k))
	r.vLen = uint8(len(v))
	r.KeyValue[0] = r.kLen
	r.KeyValue[1] = r.vLen
	copy(r.KeyValue[2:33], k)
	copy(r.KeyValue[33:], v)
	if ok{
		db.mi.put(k, inx)
		return db.write(r, int64(inx))
	}else{
		fi,_ := db.diskFile.Stat()
		db.mi.put(k, fi.Size())
		return db.writeEnd(r)
	}
}

func (db *DB) Read(k []byte) []byte{
	inx,ok := db.mi.get(k)
	if !ok{
		return nil
	}
	buf := [290]byte{}
	db.diskFile.ReadAt(buf[:], int64(inx))
	return buf[33 : 33+buf[1]]
}

func (db *DB) Close(){
	db.diskFile.Sync()
	db.diskFile.Close()
}



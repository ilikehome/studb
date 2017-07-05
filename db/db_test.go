package db

import (
	"testing"
	"os"
)

const (
	dbFolder = "test"
	dbName = "1"
	dbPath = dbFolder + string(os.PathSeparator) + dbName
)

func TestBaseRW(t *testing.T)  {
	os.MkdirAll(dbFolder, 0777)
	defer os.RemoveAll(dbFolder)

	db := Open(dbPath)
	defer db.Close()

	kStr := "gggg24r"
	k := []byte(kStr)
	vStr := "lllrrrlll4"
	v := []byte(vStr)

	db.Write(k, v)
	buf,_ := db.Read(k)
	if string(buf) != vStr{
		t.Fatal("put get failed.")
	}

	vNewStr := "tttt"
	vNew := []byte(vNewStr)

	db.Write(k,vNew)
	buf,_ = db.Read(k)
	if string(buf) != vNewStr{
		t.Fatal("update failed.")
	}

	kNotExist := []byte("cccccc")

	buf,_ = db.Read(kNotExist)
	if buf != nil{
		t.Fatal("get not exist key failed.")
	}
}



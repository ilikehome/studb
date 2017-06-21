package db

import (
	"testing"
	"fmt"
)

func TestingLogger(t *testing.T)  {
	db := Open("d:\\shdb1\\1.txt")
	defer db.Close()

	db.Write([]byte("gggg24r"),[]byte("lllrrrlll4"))
	fmt.Println(string(db.Read([]byte("gggg24r"))))
}



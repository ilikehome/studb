package main

import (
	"github.com/ilikehome/studb/studb"
	"fmt"
)


func main() {
	db := studb.Load("c:\\shdb1\\1.txt")
	defer db.Close()

	db.Write([]byte("gggg24r"),[]byte("lllrrrlll4"))
	fmt.Println(string(db.Read([]byte("gggg24r"))))
}

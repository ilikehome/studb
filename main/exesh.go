package main

import (
	"github.com/ilikehome/studb"
	"fmt"
)


func main() {
	db := studb.Load("d:\\shdb1\\1.txt")
	defer db.Close()

	db.Write([]byte("gggg24r"),[]byte("lllrrrlll4"))
	fmt.Println(string(db.Read([]byte("gggg24r"))))
}

package main

import (
	"github.com/ilikehome/studb/shdb"
	"fmt"
)


func main() {
	db := shdb.Load("/tmp/s.txt")
	defer db.Close()

	db.Write([]byte("gggg24r"),[]byte("lllrrrlll4"))
	fmt.Println(string(db.Read([]byte("gggg24r"))))
}

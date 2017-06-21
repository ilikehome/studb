package main

import (
	"github.com/ilikehome/studb/db"
	"fmt"
)


func main() {
	db := db.Open("c:\\shdb1\\1.txt")
	defer db.Close()

	db.Write([]byte("gggg24r"),[]byte("lllrrrlll4"))
	buf,_ := db.Read([]byte("gggg24r"))
	fmt.Println(string(buf))
}

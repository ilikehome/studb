package studb

import (
	"github.com/ilikehome/studb/index"
	"fmt"
)


func main() {
	db := index.Load("d:\\shdb1\\1.txt")
	defer db.Close()

	db.Write([]byte("gggg24r"),[]byte("lllrrrlll4"))
	fmt.Println(string(db.Read([]byte("gggg24r"))))
}

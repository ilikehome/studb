package index

import (
	"testing"
	"encoding/hex"
	"os"
	"github.com/ilikehome/studb/db/constant"
)

const (
	indexInDiskPath = "/1.inx.txt"
)

func TestBase(t *testing.T) {
	f,_ := os.OpenFile(indexInDiskPath, os.O_RDWR|os.O_CREATE, 0666);
	iid := openDiskIndex(f)
	defer iid.close()

	iid.append(1, constant.OP_PUT, []byte{'a', 'b', 'c'}, 64)
	iid.append(2, constant.OP_PUT, []byte{'a', 'b', 'w'},164)
	iid.append(3, constant.OP_PUT, []byte{'a', 'b', 'w'},264)
	iid.append(4, constant.OP_PUT, []byte{'a', 'b', 'w'},364)
	iid.append(5, constant.OP_PUT, []byte{'a', 'y', 'c'},464)

	m := iid.readToMem()
	if m[hex.EncodeToString([]byte{'a','y','c'})] != 464{
		t.Fatal("Load index from disk error.")
	}
}


package index

import (
	"testing"
	"encoding/hex"
)

const (
	indexInDiskPath = "d:\\shdb1\\1.inx.txt"
)

func TestBase(t *testing.T) {
	iid := Open(indexInDiskPath)
	defer iid.close()

	iid.append(64,1 ,[]byte{'a', 'b', 'c'})
	iid.append(164, 2, []byte{'a', 'b', 'w'})
	iid.append(264, 3, []byte{'a', 'b', 'w'})
	iid.append(364, 4, []byte{'a', 'b', 'w'})
	iid.append(464, 5, []byte{'a', 'y', 'c'})

	m := iid.readToMem()
	if m[hex.EncodeToString([]byte{'a','y','c'})] != 464{
		t.Fatal("Load index from disk error.")
	}
}


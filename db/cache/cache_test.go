package cache

import (
	"testing"
)


func TestBase(t *testing.T) {
	c := CreateCache(10)
	c.Put([3]byte{1,2,3}, "i")
	v := c.Get([3]byte{1,2,3})
	if v != "i" {
		t.Fatal("Get fail.")
	}

	c.Del([3]byte{1,2,3})
	v = c.Get([3]byte{1,2,3})
	if v != nil {
		t.Fatal("Del fail.")
	}
}


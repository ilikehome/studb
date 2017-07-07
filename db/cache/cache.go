package cache

import (
	"encoding/hex"
)

type Cache struct{
	m map[interface{}]interface{} //change to sync.Map in go1.9
}

func CreateCache(size int) *Cache{
	c := new(Cache)
	c.m = make(map[interface{}]interface{})
	return c
}

func(c *Cache)Put(k interface{}, v interface{}){
	value, ok := k.([]byte)
	if ok {
		c.m[hex.EncodeToString(value)] = v
	}else{
		c.m[k] = v
	}
}

func(c *Cache)Del(k interface{}){
	value, ok := k.([]byte)
	if ok {
		delete(c.m, hex.EncodeToString(value))
	}else{
		delete(c.m, k)
	}
}

func(c *Cache)Get(k interface{}) interface{}{
	value, ok := k.([]byte)
	if ok {
		return c.m[hex.EncodeToString(value)]
	}else{
		return c.m[k]
	}
}
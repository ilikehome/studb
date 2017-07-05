package cache

import "golang.org/x/tools/go/gcimporter15/testdata"

type Cache struct{
	m map[interface{}]interface{} //change to sync.Map in go1.9
}

func CreateCache(size int) *Cache{
	c := new(Cache)
	c.m = make(map[interface{}]interface{})
	return c
}

func(c *Cache)Put(k interface{}, v interface{}){
	value, ok := k.([]interface{})
	if ok {
		va := [len(value)]interface{}{}
		copy(va[:], value)
	}
	c.m[k] = v
}

func(c *Cache)Del(k interface{}){
	delete(c.m, k)
}

func(c *Cache)Get(k interface{}) interface{}{
	return c.m[k]
}
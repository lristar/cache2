package module

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type CacheLru struct {
	sync.RWMutex
	mx *sync.Mutex
	// The table's name.
	name string
	// table type
	types CacheTableType
	// All cached items.
	items map[string]*CacheItem
	// cap
	cap int
	reloadTime time.Duration
	loadData func(key string,arg ...interface{})(*CacheItem,error)
	// The logger used for this table.
	logger *log.Logger
	checkList list.List
}

func (c *CacheLru)Value(key string)interface{}{
	return nil
}

func (c *CacheLru)Delete(){

}
package module

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type CacheTableType int

var tableMap =map[string]*CacheSimple{}

type ICacheTable interface {
	Value(key string)interface{}
	Delete()
}

type CacheSimple struct {
	sync.RWMutex
	mx *sync.Mutex
	// The table's name.
	name string
	// table type
	types CacheTableType
	// All cached items.
	items map[string]*CacheItem

	reloadTime time.Duration

	loadData func(key string,arg ...interface{})(*CacheItem,error)
	// The logger used for this table.
	logger *log.Logger
	checkList list.List
}

func NewDefaultCacheTable()*CacheSimple{
	if v,ok:=tableMap[DefaultName];ok{
		return v
	}

	t :=&CacheSimple{
		mx: &sync.Mutex{},
		name:DefaultName,
		items:make(map[string]*CacheItem,0),
		reloadTime: 15*time.Second,
		loadData:nil,
		logger: log.Default(),
	}
	tableMap[DefaultName] = t
	go t.watch()
	return t
}

func (c *CacheSimple)SetReloadTime(reload time.Duration)*CacheSimple{
	c.reloadTime = reload
	return c
}

func (c *CacheSimple)SetLoadDataFun(loadData func(key string,arg ...interface{})(*CacheItem,error))*CacheSimple{
	c.loadData = loadData
	return c
}

func (c *CacheSimple)Value(key string)interface{}{
	c.RLock()
	if v,ok:=c.items[key];ok{
		return v
	}
	c.RUnlock()
	Item,err:=c.loadData(key)
	if err!=nil{
		return nil
	}
	if Item == nil{
		return nil
	}
	c.addItem(key,Item)
	return Item.Value
}

func (c *CacheSimple)Delete(){

}

func (c *CacheSimple)addItem(key string,cacheItem *CacheItem){
	cacheItem.ExpiredTime=time.Now().Add(c.reloadTime)
	c.Lock()
	c.items[key] = cacheItem
	c.Unlock()
	c.addCheckList(cacheItem,nil)
}

func (c *CacheSimple)addCheckList(cacheItem *CacheItem,element *list.Element){
	c.mx.Lock()
	defer c.mx.Unlock()
	if element!=nil{
		c.checkList.MoveToFront(element)
		return
	}
	c.checkList.PushFront(cacheItem)
}

func (c *CacheSimple)watch(){

	c.logger.Println("watching .........")

	for c.checkList.Back() == nil || c.loadData ==nil{

	}
	c.logger.Println("start .........")
	for c.checkList.Back()!=nil{
		ele :=c.checkList.Back()
		if v,ok :=ele.Value.(*CacheItem);ok{
			for v.ExpiredTime.UnixNano() > time.Now().UnixNano(){
			}
			newData ,err := c.loadData(v.Key)
			if err == nil{
				// 将新数据更新
				v.Value = newData
				v.ExpiredTime=time.Now().Add(c.reloadTime)
			}
			// 放到队列的头部
			c.addCheckList(nil,ele)
		}
	}
}


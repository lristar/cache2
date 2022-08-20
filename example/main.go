package main

import (
	"fmt"
	"myproject/cache2/inter/module"
	"time"
)

func main() {
	f :=func(key string,arg ...interface{})(*module.CacheItem,error){
		v :=time.Now().Format("2006-01-02 15:04:05 MST Mon")
		if time.Now().Second()==15{
			return nil,fmt.Errorf("失败")
		}
		return &module.CacheItem{
			Key:   key,
			Value: v,
		},nil
	}
	module.NewDefaultCacheTable(f)
	a := module.NewDefaultCacheTable(f)
	fmt.Println(a.Value("aaa"))
	time.Sleep(time.Second*5)
	fmt.Println(a.Value("bbb"))
	for {
		time.Sleep(time.Second*10)
		fmt.Println(a.Value("aaa"))
		fmt.Println(a.Value("bbb"))
	}
}

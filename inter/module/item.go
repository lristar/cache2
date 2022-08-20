package module

import "time"

type  CacheItem struct {
	Key string
	Value interface{}
	ExpiredTime time.Time
}
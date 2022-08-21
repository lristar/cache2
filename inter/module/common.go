package module


// default table name
const DefaultName  = "default"

type CacheType int
const (
	LRUCACHE = iota
	SIMPLECACHE
)
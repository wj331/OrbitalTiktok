package routes

import "sync"

var (
	Pools = map[string]*sync.Pool{}
)

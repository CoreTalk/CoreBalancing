package main

import (
	"sync"
)

var Servers []Server

var pool = &sync.Pool{
	New: func() interface{} {
		return new(Server)
	},
}

type Server struct {
	Address string
	Load    int64
}

func clear_servers() {
	for _, s := range Servers {
		pool.Put(s)
	}
}

func main() {

}

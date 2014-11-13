// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

var (
	lock                 = new(sync.Mutex)
	min_load_index int64 = 0
	Servers        servers
)

type servers []*Server

func (s servers) Len() int {
	return len(s)
}

func (s servers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s servers) Less(i, j int) bool {
	return s[i].Load < s[j].Load
}

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
	// Read the conf.
	if err := pare_config(); err != nil {
		log.Fatal(err)
	}
	// Connect the etcd.
	if err := link_etcd(); err != nil {
		log.Fatal(err)
	}
	if err := get_machines(); err != nil {
		log.Fatal(err)
	}
	// Do the flush loop.
	go flush()
	go http_listen()
	HandleSignal(InitSignal())
}

func http_listen() {
	http.HandleFunc("/get_server", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		// Return the min load server.
		w.Write([]byte(get_min_load_server()))
	})
	http.ListenAndServe(Conf.Listen_addr, nil)
}

func get_min_load_server() string {
	lock.Lock()
	s := Servers[min_load_index]
	s.Load++
	if Servers[min_load_index].Load > Servers[min_load_index+1].Load {
		min_load_index++
	}
	lock.Unlock()
	return s.Address
}

func flush() {
	C := time.After(time.Duration(Conf.Interval))
	for {
		select {
		case <-C:
			get_machines()
			sort_servers()
		}
	}
}

func sort_servers() {
	lock.Lock()
	sort.Sort(Servers)
	lock.Unlock()
}

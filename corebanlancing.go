// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var (
	lock                 = new(sync.Mutex)
	min_load_index int64 = 0
	Servers        servers
	temp_servers   servers
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

func clear_servers(ss *servers) {
	for _, s := range *ss {
		pool.Put(s)
	}
	*ss = (*ss)[:0]
}

func set_list(nodes []*etcd.Node, ss servers) {
	for _, n := range nodes {
		s := pool.Get().(*Server)
		s.Address = n.Key
		s.Load, _ = strconv.ParseInt(n.Value, 10, 64)
		ss = append(ss, s)
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
	if err := get_machines(Servers); err != nil {
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
	if s.Load > Servers[min_load_index+1].Load {
		if min_load_index < int64(len(Servers)) {
			min_load_index++
		}
	}
	lock.Unlock()
	return s.Address
}

func flush() {
	C := time.After(time.Duration(Conf.Interval))
	for {
		select {
		case <-C:
			clear_servers(&temp_servers)
			get_machines(temp_servers)
			sort_servers(temp_servers)
			lock.Lock()
			Servers = temp_servers
			lock.Unlock()
			temp_servers = make(servers, 0, len(Servers))
		}
	}
}

func sort_servers(ss servers) {
	sort.Sort(ss)
}

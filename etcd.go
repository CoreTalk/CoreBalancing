// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var Client *etcd.Client
var stop_watch = make(chan bool)

func link_etcd() error {
	Client = etcd.NewClient(Conf.Etcd_machines)
	if _, err := Client.CreateDir(Conf.Node_name, 0); err != nil {
		return
	}
	_, err := Client.Set("/core_bl/"+Conf.Listen_addr, "1", strconv.FormatInt(Conf.Heart_beat_time, 10))
	return err
}

func heart_beat() chan error {
	c := time.After(Conf.Heart_beat_time / 2)
	ch := make(chan error)
	for {
		select {
		case <-c:
			_, err := Client.Update("/core_bl/"+Conf.Listen_addr, "1", strconv.FormatInt(Conf.Heart_beat_time, 10))
			if err != nil {
				ch <- err
			}
		}
	}
	return ch
}

func get_machines() error {
	resp, err := Client.Get(Conf.Node_name, false, false)
	if err != nil {
		return err
	}
	set_list(resp.Node.Nodes)
	return nil
}

func set_list(nodes []*etcd.Node) {
	for _, n := range nodes {
		s := pool.Get().(*Server)
		s.Address = n.Key
		s.Load, _ = strconv.ParseInt(n.Value, 10, 64)
	}
}

// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"

	"github.com/coreos/go-etcd/etcd"
)

var Client *etcd.Client
var stop_watch = make(chan bool)

func link_etcd() error {
	Client = etcd.NewClient(Conf.Etcd_machines)
	_, err := Client.CreateDir(Conf.Node_name, 0)
	return err
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

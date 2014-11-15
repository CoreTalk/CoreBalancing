// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var Client *etcd.Client
var stop_watch = make(chan bool)

func link_etcd() error {
	Client = etcd.NewClient(Conf.Etcd_machines)
	if _, err := Client.CreateDir(Conf.Node_name, 0); err != nil {
		return err
	}
	_, err := Client.Set("/core_bl/"+Conf.Listen_addr, "running", Conf.Heart_beat_time)
	return err
}

func heart_beat() chan error {
	c := time.After(time.Duration(Conf.Heart_beat_time / 2))
	ch := make(chan error)
	go func() {
		for {
			select {
			case <-c:
				_, err := Client.Update("/core_bl/"+Conf.Listen_addr, "running", Conf.Heart_beat_time)
				if err != nil {
					ch <- err
				}
			}
		}
	}()
	return ch
}

func get_machines(ss servers) error {
	resp, err := Client.Get(Conf.Node_name, false, false)
	if err != nil {
		return err
	}
	set_list(resp.Node.Nodes, ss)
	return nil
}

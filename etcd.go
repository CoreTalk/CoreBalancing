// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var Client *etcd.Client
var stop_watch = make(chan bool)

func link_etcd() error {
	Client = etcd.NewClient(Conf.Balancer.Etcd_machines)
	if _, err := Client.CreateDir(Conf.Balancer.Node_name, 0); err != nil {
		if err != nil {
			if e, ok := err.(*etcd.EtcdError); ok {
				if e.ErrorCode != 105 {
					return err
				}
				log.Println("Key has been exist.")
			} else {
				return err
			}
		}
	}
	_, err := Client.Set("/core_bl/"+Conf.Balancer.Listen_addr, "running", Conf.Balancer.Heart_beat_time)
	return err
}

func heart_beat() chan error {
	t := time.NewTicker(time.Duration(Conf.Balancer.Heart_beat_time/2) * time.Second)
	ch := make(chan error)
	go func() {
		for {
			select {
			case <-t.C:
				_, err := Client.Update("/core_bl/"+Conf.Balancer.Listen_addr, "running", Conf.Balancer.Heart_beat_time)
				if err != nil {
					log.Printf("Etcd hb error:%v", err)
					ch <- err
				}
			}
		}
	}()
	return ch
}

func get_machines(ss *servers) error {
	resp, err := Client.Get(Conf.Balancer.Node_name, false, false)
	if err != nil {
		return err
	}
	set_list(resp.Node.Nodes, ss)
	return nil
}

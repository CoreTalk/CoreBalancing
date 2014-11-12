package main

import (
	"time"

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

func flush() {
	C := time.After(time.Duration(Conf.Interval))
	for {
		select {
		case <-C:
			get_machines()
		}
	}
}

func set_list(nodes []etcd.Node) {
	for _, n := range nodes {
		s := pool.Get().(*Server)
		s.Address = n.Key
		s.Load = n.Value
	}
}

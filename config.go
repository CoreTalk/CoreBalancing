// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/gcfg"
)

var Conf = new(Config)

type Config struct {
	Balancer
}

type Balancer struct {
	// [scheme]://[ip]:[port]
	Etcd_machines []string
	// Interval. Sort the servers. < comet flush time.
	Interval int64
	// The etcd node's name.
	Node_name string

	// [IP]:[prot]
	Listen_addr string

	// Heart beat time.
	Heart_beat_time uint64
}

func pare_config() error {
	return gcfg.ReadFileInto(Conf, "user.conf")
}

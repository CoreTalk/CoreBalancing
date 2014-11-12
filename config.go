// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/gcfg"
)

var Conf *Config

type Config struct {
	// [scheme]://[ip]:[port]
	Etcd_machines []string
	// Interval. Sort the servers.
	Interval int64
	// The etcd node's name.
	Node_name string

	// [IP]:[prot]
	Listen_addr string
}

func pare_config() error {
	return gcfg.ReadFileInto(Conf, "conf.user.conf")
}

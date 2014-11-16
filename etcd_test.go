// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"

	"github.com/coreos/go-etcd/etcd"
)

func TestDir(t *testing.T) {
	c := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	resp, err := c.Get("/core_bl", false, false)
	if err != nil {
		t.Error(err)
	}
	if !resp.Node.Dir {
		t.Fatal("not dir.")
	}
	resp, err = c.Get("/comets", false, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Len is", resp.Node.Nodes.Len())
	for _, n := range resp.Node.Nodes {
		fmt.Println("node is", n.Key, n.Value)
	}
}

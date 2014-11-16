// Copyright © 2014 CoreTalk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

func TestParseconfig(t *testing.T) {
	if err := pare_config(); err != nil {
		t.Error(err)
	}
	fmt.Println(Conf)
}

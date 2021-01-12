// SPDX-License-Identifier: MIT
//
// Copyright © 2021 Kent Gibson <warthog618@gmail.com>.

package blinkt_test

import (
	"testing"

	"github.com/warthog618/blinkt"
)

func BenchmarkShow(b *testing.B) {
	bl := blinkt.New()
	for i := 0; i < b.N; i++ {
		bl.Show()
	}
	bl.Close()
}

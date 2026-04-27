/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package cmodule

import (
	"testing"
)

func TestEnsureBuffer(t *testing.T) {
	var seq Sequence
	err := seq.EnsureCap(3)
	if err != nil {
		t.Error("error not nil:", err.Error())
	} else if len(seq.buffer) == 0 {
		t.Error("buffer length = 0")
	} else {
		bufferPrevLen := len(seq.buffer)
		err = seq.EnsureCap(5)
		if err != nil {
			t.Error("error not nil:", err.Error())
		} else if len(seq.buffer) == 0 {
			t.Error("buffer length = 0")
		} else if bufferPrevLen == len(seq.buffer) {
			t.Error("buffer length unchanged")
		} else {
			err = seq.Close()
			if err != nil {
				t.Error("error not nil:", err.Error())
			} else if len(seq.buffer) > 0 {
				t.Error("buffer length > 0")
			}
		}
	}
}

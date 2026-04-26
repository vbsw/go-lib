/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package cmodule

import (
	"testing"
	"unsafe"
)

func TestEnsureBuffer(t *testing.T) {
	var seq Sequence
	seq.Modules = []Module{nil, nil, nil}
	err := seq.ensureBuffer()
	if err != nil {
		t.Error("error not nil:", err.Error())
	} else if len(seq.buffer) == 0 {
		t.Error("buffer length = 0")
	} else {
		bufferLenPrev := len(seq.buffer)
		seq.Modules = append(seq.Modules, nil, nil)
		seq.buffer = unsafe.Slice(&seq.buffer[0], 5)
		err = seq.ensureBuffer()
		if err != nil {
			t.Error("error not nil:", err.Error())
		} else if len(seq.buffer) == 0 {
			t.Error("buffer length = 0")
		} else if bufferLenPrev == len(seq.buffer) {
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

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

func TestNew(t *testing.T) {
	seq := NewSequence(4)
	if len(seq) == 0 {
		t.Error("C memory not initialized")
	} else if len(seq) != 4*SequenceChunks {
		t.Error("wrong total length:", len(seq))
	} else if seq.Len() != 4 {
		t.Error("wrong length:", seq.Len())
	} else {
		seq = seq.Release()
		if len(seq) != 0 {
			t.Error("wrong length:", len(seq))
		}
	}
}

func TestRemove0(t *testing.T) {
	seq := NewSequence(4)
	init4Test(seq)
	seq = seq.Remove(0)
	if len(seq) == 0 {
		t.Error("length is 0")
	} else if len(seq) != 3*SequenceChunks {
		t.Error("wrong total length:", len(seq))
	} else if seq.Len() != 3 {
		t.Error("wrong length:", seq.Len())
	} else if seq[0] != 1 {
		t.Error("wrong function pointer:", int(seq[0]))
	} else if seq[3] != 5 {
		t.Error("wrong data pointer:", int(seq[3]))
	} else if seq[6] != 9 {
		t.Error("wrong extra1:", int(seq[6]))
	} else if seq[9] != 13 {
		t.Error("wrong extra2:", int(seq[9]))
	} else {
		seq = seq.Release()
	}
}

func TestRemove1(t *testing.T) {
	seq := NewSequence(4)
	init4Test(seq)
	seq = seq.Remove(1)
	if len(seq) == 0 {
		t.Error("length is 0")
	} else if len(seq) != 3*SequenceChunks {
		t.Error("wrong total length:", len(seq))
	} else if seq.Len() != 3 {
		t.Error("wrong length:", seq.Len())
	} else if seq[1] != 2 {
		t.Error("wrong function pointer:", int(seq[1]))
	} else if seq[4] != 6 {
		t.Error("wrong data pointer:", int(seq[4]))
	} else if seq[7] != 10 {
		t.Error("wrong extra1:", int(seq[7]))
	} else if seq[10] != 14 {
		t.Error("wrong extra2:", int(seq[10]))
	} else {
		seq = seq.Release()
	}
}

func TestRemove3(t *testing.T) {
	seq := NewSequence(4)
	init4Test(seq)
	seq = seq.Remove(3)
	if len(seq) == 0 {
		t.Error("length is 0")
	} else if len(seq) != 3*SequenceChunks {
		t.Error("wrong total length:", len(seq))
	} else if seq.Len() != 3 {
		t.Error("wrong length:", seq.Len())
	} else if seq[3] != 4 {
		t.Error("wrong function pointer:", int(seq[3]))
	} else if seq[6] != 8 {
		t.Error("wrong data pointer:", int(seq[6]))
	} else if seq[9] != 12 {
		t.Error("wrong extra1:", int(seq[9]))
	} else if seq[11] != 14 {
		t.Error("wrong extra2:", int(seq[11]))
	} else {
		seq = seq.Release()
	}
}

func TestRemove13(t *testing.T) {
	seq := NewSequence(4)
	init4Test(seq)
	seq = seq.Remove(1, 3)
	if len(seq) == 0 {
		t.Error("length is 0")
	} else if len(seq) != 2*SequenceChunks {
		t.Error("wrong total length:", len(seq))
	} else if seq.Len() != 2 {
		t.Error("wrong length:", seq.Len())
	} else if seq[1] != 2 {
		t.Error("wrong function pointer:", int(seq[1]))
	} else if seq[3] != 6 {
		t.Error("wrong data pointer:", int(seq[3]))
	} else if seq[5] != 10 {
		t.Error("wrong extra1:", int(seq[5]))
	} else if seq[7] != 14 {
		t.Error("wrong extra2:", int(seq[7]))
	} else {
		seq = seq.Release()
	}
}

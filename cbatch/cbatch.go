/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package cbatch provides C batch processing.
package cbatch

// #include <stdint.h>
// #include "cbatch.h"
import "C"
import (
	"strconv"
	"unsafe"
)

const (
	Init Step = iota
	Process
	Destroy
)

const (
	maxInt         = int((^uint(0)) >> 1)
	maxInt32       = int32((^uint32(0)) >> 1)
	SequenceChunks = 4
	MaxSequenceLen = min(uint64(C.SIZE_MAX)/SequenceChunks, uint64(maxInt/SequenceChunks))
)

// Step represents a step in a task.
type Step int

// Error is returned by Run.
type Error struct {
	Num1  int64
	Num2  int64
	Str   string
	Index int
}

// Sequence holds pointers to functions and data
// This is an array allocated in C and must be released when no more needed.
type Sequence []unsafe.Pointer

// Task provides abstraction to C functions and its data.
type Task interface {
	// CData returns pointer to C function and C data.
	CData(step Step) (unsafe.Pointer, unsafe.Pointer)
	AsError(num1, num2 int64, str string) error
	SetCData(data unsafe.Pointer)
}

// NewSequence returns a new instance of Sequence.
// Sequence is an array allocated in C and must be released when no more needed.
func NewSequence(length int) Sequence {
	if length > 0 && uint64(length) <= MaxSequenceLen {
		var dataC *unsafe.Pointer
		totalLength := length * SequenceChunks
		C.cbatch_alloc(&dataC, C.size_t(totalLength))
		if dataC != nil {
			return unsafe.Slice(dataC, totalLength)
		}
		return nil
	}
	panic("sequence length unsupported")
}

// Remove removes elements from tasks array and returns it.
// Indices must be in ascending order.
func Remove(tasks []Task, indices ...int) []Task {
	if len(indices) > 0 {
		var gap, gapFrom, gapTo int
		for _, index := range indices {
			if gapFrom == gapTo {
				gapFrom, gapTo = index, index+1
			} else if gapTo == index {
				gapTo++
			} else if gapTo < index {
				copy(tasks[gapFrom-gap:], tasks[gapTo:index])
				gap += (gapTo - gapFrom)
				gapFrom, gapTo = index, index+1
			} else {
				panic("wrong indices order")
			}
		}
		if gapTo < len(tasks) {
			copy(tasks[gapFrom-gap:], tasks[gapTo:])
		}
		tasks = tasks[:len(tasks)-len(indices)]
	}
	return tasks
}

// Disable sets functions to be skipped in Run. Applies to all when indices empty.
// To enable entries back again use Setup().
func (seq Sequence) Disable(indices ...int) {
	length := seq.Len()
	if len(indices) > 0 {
		for _, index := range indices {
			seq[index], seq[index+length] = nil, nil
		}
	} else {
		for i := 0; i < length; i++ {
			seq[i], seq[i+length] = nil, nil
		}
	}
}

// Len returns number of tasks in Sequence.
func (seq Sequence) Len() int {
	return len(seq) / SequenceChunks
}

// Run processes C data in Sequence.
func (seq Sequence) Run(passes int) *Error {
	length := seq.Len()
	if passes > 0 && length > 0 {
		if uint64(passes) <= uint64(maxInt32) {
			var params C.cbatch_proc_params_t
			params.data = &seq[0]
			params.length = C.size_t(length)
			params.passes = C.int32_t(passes)
			C.cbatch_proc(&params)
			if params.err1 != 0 {
				err := new(Error)
				err.Num1 = int64(params.err1)
				err.Num2 = int64(params.err2)
				err.Index = int(params.err_idx)
				if params.err_str != nil {
					err.Str = C.GoString(params.err_str)
				}
				return err
			}
			return nil
		}
		panic("passes count unsupported")
	}
	return nil
}

// RunInit is abbreviation for Setup(Init), Run(passes),
// Sync(tasks) and AsError(tasks)
func (seq Sequence) RunInit(tasks []Task, passes int) error {
	length := seq.Len()
	for i := 0; i < length && i < len(tasks); i++ {
		seq[i], seq[i+length] = tasks[i].CData(Init)
	}
	err := seq.Run(passes)
	if err == nil {
		for i := 0; i < length && i < len(tasks); i++ {
			tasks[i].SetCData(seq[i+length])
		}
		return nil
	}
	return err.AsError(tasks)
}

// RunProcess is abbreviation for Setup(Process), Run(passes),
// Sync(tasks) and AsError(tasks)
func (seq Sequence) RunProcess(tasks []Task, passes int) error {
	length := seq.Len()
	for i := 0; i < length && i < len(tasks); i++ {
		seq[i], seq[i+length] = tasks[i].CData(Process)
	}
	err := seq.Run(passes)
	if err == nil {
		for i := 0; i < length && i < len(tasks); i++ {
			tasks[i].SetCData(seq[i+length])
		}
		return nil
	}
	return err.AsError(tasks)
}

// RunDestroy is abbreviation for Setup(Destroy), Run(passes),
// Sync(tasks) and AsError(tasks)
func (seq Sequence) RunDestroy(tasks []Task, passes int) error {
	length := seq.Len()
	for i := 0; i < length && i < len(tasks); i++ {
		seq[i], seq[i+length] = tasks[i].CData(Destroy)
	}
	err := seq.Run(passes)
	if err == nil {
		for i := 0; i < length && i < len(tasks); i++ {
			tasks[i].SetCData(seq[i+length])
		}
		return nil
	}
	return err.AsError(tasks)
}

// Release releases C memory. Returns always nil.
func (seq Sequence) Release() Sequence {
	if len(seq) > 0 {
		C.cbatch_free(&seq[0])
	}
	return nil
}

// Remove removes elements from Sequence and returns it.
// Indices must be in ascending order, in interval [0,seq.Len()) and must not remove everything.
// (There is no function, yet, to insert them back again.)
func (seq Sequence) Remove(indices ...int) Sequence {
	if len(indices) > 0 {
		if length := seq.Len(); length > len(indices) {
			var gap, gapFrom, gapTo int
			// entries are removed from each chunk and moved to front to close the gap
			for i := 0; i < SequenceChunks; i++ {
				for _, index := range indices {
					offIdx := length*i + index
					if gapFrom == gapTo {
						gapFrom, gapTo = offIdx, offIdx+1
					} else if gapTo == offIdx {
						gapTo++
					} else if gapTo < offIdx {
						copy(seq[gapFrom-gap:], seq[gapTo:offIdx])
						gap += (gapTo - gapFrom)
						gapFrom, gapTo = offIdx, offIdx+1
					} else {
						panic("wrong indices order")
					}
				}
			}
			// move rest of entries (of last chunk) to close the gap
			if gapTo < len(seq) {
				copy(seq[gapFrom-gap:], seq[gapTo:])
			}
			// adjust length
			seq = seq[:len(seq)-len(indices)*SequenceChunks]
		} else {
			panic("wrong indices length")
		}
	}
	return seq
}

// Setup sets functions and data for Run. Applies to all when indices empty.
func (seq Sequence) Setup(step Step, tasks []Task, indices ...int) {
	length := seq.Len()
	if len(indices) > 0 {
		for _, index := range indices {
			seq[index], seq[index+length] = tasks[index].CData(step)
		}
	} else {
		for i := 0; i < length && i < len(tasks); i++ {
			seq[i], seq[i+length] = tasks[i].CData(step)
		}
	}
}

// Sync writes C data to tasks. Applies to all when indices empty.
func (seq Sequence) Sync(tasks []Task, indices ...int) {
	length := seq.Len()
	if len(indices) > 0 {
		for _, index := range indices {
			tasks[index].SetCData(seq[index+length])
		}
	} else {
		for i := 0; i < length && i < len(tasks); i++ {
			tasks[i].SetCData(seq[i+length])
		}
	}
}

// AsError converts cbatch's Error to Go error and returns it.
func (batchErr *Error) AsError(tasks []Task) error {
	err := tasks[batchErr.Index].AsError(batchErr.Num1, batchErr.Num2, batchErr.Str)
	if err != nil {
		return err
	}
	return batchErr
}

func (batchErr *Error) Error() string {
	var str string
	if batchErr.Num1 < 1000000 {
		str = "out of memory"
	} else {
		str = "unknown"
	}
	str = str + " (" + strconv.FormatInt(batchErr.Num1, 10)
	if batchErr.Num2 == 0 {
		str = str + ")"
	} else {
		str = str + ", " + strconv.FormatInt(batchErr.Num2, 10) + ")"
	}
	if len(batchErr.Str) > 0 {
		str = str + "; " + batchErr.Str
	}
	return str
}

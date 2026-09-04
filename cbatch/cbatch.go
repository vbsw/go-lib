/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package cbatch provides C batch processing.
package cbatch

// #include <limits.h>
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
	maxIntFast32   = C.int_fast32_t((^C.uint_fast32_t(0)) >> 1)
	SequenceChunks = 4
	MaxSequenceLen = (min(uint64(maxIntFast32), uint64(C.SIZE_MAX)) / SequenceChunks) / uint64(unsafe.Sizeof(unsafe.Pointer(nil)))
)

// Step represents a step in a task.
type Step int

// Error is returned by Run.
type Error struct {
	Str      string
	TaskName string
	Num1     int
	Num2     int
	Index    int
}

// Sequence holds pointers to functions and data
// This is an array allocated in C and must be released when no more needed.
type Sequence []unsafe.Pointer

// Task provides abstraction to C functions and its data.
type Task interface {
	// CData returns pointer to C function and C data1, data2.
	CData(step Step) (unsafe.Pointer, unsafe.Pointer, unsafe.Pointer)
	AsError(num1, num2 int, str string) error
	SetCData(step Step, data1, data2 unsafe.Pointer)
}

// NewSequence returns a new instance of Sequence.
// Sequence is an array allocated in C and must be released when no more needed.
func NewSequence(length int) Sequence {
	if length > 0 {
		if uint64(length) <= MaxSequenceLen {
			var dataC *unsafe.Pointer
			dataLen := length * SequenceChunks
			C.aca_batch_alloc(&dataC, C.int_fast32_t(dataLen))
			if dataC != nil {
				return unsafe.Slice(dataC, dataLen)
			}
			return nil
		}
		panic("sequence length overflow")
	}
	panic("sequence length underflow")
}

// Remove removes elements from tasks array and returns this modified tasks array.
//
// Precondition:
//
//	len(tasks) >= len(indices) && indices[i] < indices[i+1] && 0 <= indices[i] < len(tasks)
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
				panic("indices not in ascending order")
			}
		}
		if gapTo < len(tasks) {
			copy(tasks[gapFrom-gap:], tasks[gapTo:])
		}
		tasks = tasks[:len(tasks)-len(indices)]
	}
	return tasks
}

// Disable sets functions to nil so they are skipped in Run.
// Applies to all when indices empty.
// To enable entries back again use Setup().
//
// Precondition:
//
//	seq.Len() >= len(indices) && 0 <= indices[i] < seq.Len()
func (seq Sequence) Disable(indices ...int) {
	seqLen := seq.Len()
	if len(indices) > 0 {
		for _, index := range indices {
			seq[index], seq[seqLen+index], seq[seqLen*2+index] = nil, nil, nil
		}
	} else {
		for i := 0; i < seqLen; i++ {
			seq[i], seq[seqLen+i], seq[seqLen*2+i] = nil, nil, nil
		}
	}
}

// Len returns maximum number of tasks in Sequence.
func (seq Sequence) Len() int {
	return len(seq) / SequenceChunks
}

// Run processes C data in Sequence.
func (seq Sequence) Run(passes int) *Error {
	seqLen := seq.Len()
	if passes > 0 && seqLen > 0 {
		if uint64(passes) <= uint64(maxIntFast32) {
			var err1, err2, err_idx C.int_fast32_t
			var err_str *C.char
			C.aca_batch_run(&seq[0], &err1, &err2, &err_idx, &err_str, C.int_fast32_t(seqLen), C.int_fast32_t(passes))
			if err1 != 0 {
				err := new(Error)
				taskNameC := seq[seqLen*3+err.Index]
				err.Index = int(err_idx)
				err.Num1 = int(err1)
				err.Num2 = int(err2)
				if err_str != nil {
					err.Str = C.GoString(err_str)
				}
				if taskNameC != nil {
					err.TaskName = C.GoString((*C.char)(taskNameC))
				}
				return err
			}
			return nil
		}
		panic("passes count overflow")
	}
	return nil
}

// RunInit is abbreviation for Setup(Init), Run(passes),
// Sync(tasks) and AsError(tasks)
func (seq Sequence) RunInit(tasks []Task, passes int) error {
	seqLen := seq.Len()
	for i := 0; i < seqLen && i < len(tasks); i++ {
		seq[i], seq[seqLen+i], seq[seqLen*2+i] = tasks[i].CData(Init)
	}
	err := seq.Run(passes)
	if err == nil {
		for i := 0; i < seqLen && i < len(tasks); i++ {
			tasks[i].SetCData(Init, seq[seqLen+i], seq[seqLen*2+i])
		}
		return nil
	}
	return err.AsError(tasks)
}

// RunProcess is abbreviation for Setup(Process), Run(passes),
// Sync(tasks) and AsError(tasks)
func (seq Sequence) RunProcess(tasks []Task, passes int) error {
	seqLen := seq.Len()
	for i := 0; i < seqLen && i < len(tasks); i++ {
		seq[i], seq[seqLen+i], seq[seqLen*2+i] = tasks[i].CData(Process)
	}
	err := seq.Run(passes)
	if err == nil {
		for i := 0; i < seqLen && i < len(tasks); i++ {
			tasks[i].SetCData(Process, seq[seqLen+i], seq[seqLen*2+i])
		}
		return nil
	}
	return err.AsError(tasks)
}

// RunDestroy is abbreviation for Setup(Destroy), Run(passes),
// Sync(tasks) and AsError(tasks)
func (seq Sequence) RunDestroy(tasks []Task, passes int) error {
	seqLen := seq.Len()
	for i := 0; i < seqLen && i < len(tasks); i++ {
		seq[i], seq[seqLen+i], seq[seqLen*2+i] = tasks[i].CData(Destroy)
	}
	err := seq.Run(passes)
	if err == nil {
		for i := 0; i < seqLen && i < len(tasks); i++ {
			tasks[i].SetCData(Destroy, seq[seqLen+i], seq[seqLen*2+i])
		}
		return nil
	}
	return err.AsError(tasks)
}

// Release releases C memory. Returns always nil.
func (seq Sequence) Release() Sequence {
	if len(seq) > 0 {
		C.aca_batch_free(&seq[0])
	}
	return nil
}

// Remove removes elements from Sequence and returns this modified Sequence.
// (There is no function, yet, to insert them back again.)
//
// Precondition:
//
//	len(indices) == 0 || seq.Len() > len(indices) && indices[i] < indices[i+1] && 0 <= indices[i] < seq.Len()
func (seq Sequence) Remove(indices ...int) Sequence {
	if len(indices) > 0 {
		if seqLen := seq.Len(); len(indices) < seqLen {
			var gap, gapFrom, gapTo int
			// entries are removed from each chunk and moved to front to close the gap
			for i := 0; i < SequenceChunks; i++ {
				for _, index := range indices {
					offIdx := seqLen*i + index
					if gapFrom == gapTo {
						gapFrom, gapTo = offIdx, offIdx+1
					} else if gapTo == offIdx {
						gapTo++
					} else if gapTo < offIdx {
						copy(seq[gapFrom-gap:], seq[gapTo:offIdx])
						gap += (gapTo - gapFrom)
						gapFrom, gapTo = offIdx, offIdx+1
					} else {
						panic("indices not in ascending order")
					}
				}
			}
			// move rest of entries (of last chunk) to close the gap
			if gapTo < len(seq) {
				copy(seq[gapFrom-gap:], seq[gapTo:])
			}
			// adjust seqLen
			seq = seq[:len(seq)-len(indices)*SequenceChunks]
		} else {
			panic("indices length overflow")
		}
	}
	return seq
}

// Setup sets functions and data for Run. Applies to all when indices empty.
//
// Precondition:
//
//	0 <= indices[i] < seq.Len() && 0 <= indices[i] < len(tasks)
func (seq Sequence) Setup(step Step, tasks []Task, indices ...int) {
	seqLen := seq.Len()
	if len(indices) > 0 {
		for _, index := range indices {
			seq[index], seq[seqLen+index], seq[seqLen*2+index] = tasks[index].CData(step)
		}
	} else {
		for i := 0; i < seqLen && i < len(tasks); i++ {
			seq[i], seq[seqLen+i], seq[seqLen*2+i] = tasks[i].CData(step)
		}
	}
}

// Sync writes C data to tasks. Applies to all when indices empty.
//
// Precondition:
//
//	0 <= indices[i] < seq.Len() && 0 <= indices[i] < len(tasks)
func (seq Sequence) Sync(step Step, tasks []Task, indices ...int) {
	seqLen := seq.Len()
	if len(indices) > 0 {
		for _, index := range indices {
			tasks[index].SetCData(step, seq[seqLen+index], seq[seqLen*2+index])
		}
	} else {
		for i := 0; i < seqLen && i < len(tasks); i++ {
			tasks[i].SetCData(step, seq[seqLen+i], seq[seqLen*2+i])
		}
	}
}

// AsError passes cbatch's Error to task's AsError it came from and returns its result.
//
// Precondition:
//
//	0 <= batchErr.Index < len(tasks)
func (batchErr *Error) AsError(tasks []Task) error {
	err := tasks[batchErr.Index].AsError(batchErr.Num1, batchErr.Num2, batchErr.Str)
	if err != nil {
		return err
	}
	return batchErr
}

// Error converts Error data to string and returns it.
func (batchErr *Error) Error() string {
	var str string
	if batchErr.Num1 < 1000000 {
		str = "out of memory"
	} else {
		str = "unknown"
	}
	str = str + " (" + strconv.FormatInt(int64(batchErr.Num1), 10)
	if batchErr.Num2 == 0 {
		str = str + ")"
	} else {
		str = str + ", " + strconv.FormatInt(int64(batchErr.Num2), 10) + ")"
	}
	if len(batchErr.Str) > 0 {
		str = str + "; " + batchErr.Str
	}
	return str
}

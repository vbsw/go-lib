/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package cmodule provides C batch processing.
package cmodule

// #include <stdint.h>
// #include "cmodule.h"
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

const (
	Init SequenceType = iota
	Run
	Destroy
)

const (
	maxInt         = int((^uint(0)) >> 1)
	maxInt32       = int32((^uint32(0)) >> 1)
	SequenceChunks = 4
	MaxSequenceLen = min(uint64(C.SIZE_MAX)/SequenceChunks, uint64(maxInt/SequenceChunks))
)

// SequenceType represents processing type.
type SequenceType int

// Error is returned by Process.
type Error struct {
	ModuleErr int64
	SystemErr int64
	Info      string
	Index     int
}

// Sequence is a buffer allocated in C.
// It holds pointers to functions and data.
type Sequence []unsafe.Pointer

// Module provides abstraction to C functions and data.
type Module interface {
	// CProcessor returns pointer to a C function and C data.
	CProcessor(seqenceType SequenceType) (unsafe.Pointer, unsafe.Pointer)
	CToGoError(moduleErr, systemErr int64, info string) error
	SetCData(data unsafe.Pointer)
}

// NewSequence returns a new instance of Sequence.
// Memory in C is allocated.
func NewSequence(length int) Sequence {
	if length > 0 && uint64(length) <= MaxSequenceLen {
		var dataC *unsafe.Pointer
		totalLength := length * SequenceChunks
		C.cmodule_alloc(&dataC, C.size_t(totalLength))
		if dataC != nil {
			return unsafe.Slice(dataC, totalLength)
		}
		return nil
	}
	panic("sequence length unsupported")
}

// Disable sets functions to be skipped in Process. Applies to all when indices empty.
func (seq Sequence) Disable(indices ...int) {
	length := seq.Len()
	if len(indices) == 0 {
		for i := 0; i < length; i++ {
			seq[i], seq[i+length] = nil, nil
		}
	} else {
		for _, index := range indices {
			seq[index], seq[index+length] = nil, nil
		}
	}
}

// Len returns the number of modules in the Sequence.
func (seq Sequence) Len() int {
	return len(seq) / SequenceChunks
}

// Process processes C data in batch.
func (seq Sequence) Process(passes int) *Error {
	length := seq.Len()
	if passes > 0 && length > 0 {
		if uint64(passes) <= uint64(maxInt32) {
			var params C.cmodule_proc_params_t
			params.data = &seq[0]
			params.length = C.size_t(length)
			params.passes = C.int32_t(passes)
			C.cmodule_proc(&params)
			if params.err1 != 0 {
				err := new(Error)
				err.ModuleErr = int64(params.err1)
				err.SystemErr = int64(params.err2)
				err.Index = int(params.err_idx)
				if params.err_str != nil {
					err.Info = C.GoString(params.err_str)
				}
				return err
			}
			return nil
		}
		panic("passes count unsupported")
	}
	return nil
}

// ProcessInit is abbreviation for Set(Init), Process(passes),
// Sync(modules) and GoError(modules)
func (seq Sequence) ProcessInit(modules []Module, passes int) error {
	length := seq.Len()
	for i := 0; i < length && i < len(modules); i++ {
		seq[i], seq[i+length] = modules[i].CProcessor(Init)
	}
	err := seq.Process(passes)
	if err == nil {
		for i := 0; i < length && i < len(modules); i++ {
			modules[i].SetCData(seq[i+length])
		}
		return nil
	}
	return err.GoError(modules)
}

// ProcessRun is abbreviation for Set(Run), Process(passes),
// Sync(modules) and GoError(modules)
func (seq Sequence) ProcessRun(modules []Module, passes int) error {
	length := seq.Len()
	for i := 0; i < length && i < len(modules); i++ {
		seq[i], seq[i+length] = modules[i].CProcessor(Run)
	}
	err := seq.Process(passes)
	if err == nil {
		for i := 0; i < length && i < len(modules); i++ {
			modules[i].SetCData(seq[i+length])
		}
		return nil
	}
	return err.GoError(modules)
}

// ProcessDestroy is abbreviation for Set(Destroy), Process(passes),
// Sync(modules) and GoError(modules)
func (seq Sequence) ProcessDestroy(modules []Module, passes int) error {
	length := seq.Len()
	for i := 0; i < length && i < len(modules); i++ {
		seq[i], seq[i+length] = modules[i].CProcessor(Destroy)
	}
	err := seq.Process(passes)
	if err == nil {
		for i := 0; i < length && i < len(modules); i++ {
			modules[i].SetCData(seq[i+length])
		}
		return nil
	}
	return err.GoError(modules)
}

// Release releases C memory. Returns always nil.
func (seq Sequence) Release() Sequence {
	if len(seq) > 0 {
		C.cmodule_free(&seq[0])
	}
	return nil
}

// Remove removes elements from Sequence. Indices must be in ascending order
// and must not remove everything.
func (seq Sequence) Remove(indices ...int) Sequence {
	if len(indices) > 0 {
		length := seq.Len()
		if length > len(indices) {
			delta, i1, i2 := seq.remove(0, 0, 0, 0, indices)
			delta, i1, i2 = seq.remove(delta, i1, i2, length, indices)
			delta, i1, i2 = seq.remove(delta, i1, i2, length*2, indices)
			delta, i1, i2 = seq.remove(delta, i1, i2, length*3, indices)
			if i2 < len(seq) {
				copy(seq[i1-delta:], seq[i2:])
			}
			seq = seq[:len(seq)-len(indices)*SequenceChunks]
		} else {
			panic("wrong indices length")
		}
	}
	return seq
}

// Set sets functions and data for Process. Applies to all when indices empty.
func (seq Sequence) Set(seqenceType SequenceType, modules []Module, indices ...int) {
	length := seq.Len()
	if len(indices) == 0 {
		for i := 0; i < length && i < len(modules); i++ {
			seq[i], seq[i+length] = modules[i].CProcessor(seqenceType)
		}
	} else {
		for _, index := range indices {
			seq[index], seq[index+length] = modules[index].CProcessor(seqenceType)
		}
	}
}

// Sync writes C data to modules. Applies to all when indices empty.
func (seq Sequence) Sync(modules []Module, indices ...int) {
	length := seq.Len()
	if len(indices) == 0 {
		for i := 0; i < length && i < len(modules); i++ {
			modules[i].SetCData(seq[i+length])
		}
	} else {
		for _, index := range indices {
			modules[index].SetCData(seq[index+length])
		}
	}
}

func (seq Sequence) remove(delta, i1, i2, offset int, indices []int) (int, int, int) {
	for _, index := range indices {
		offIdx := offset + index
		if i1 == i2 {
			i1, i2 = offIdx, offIdx+1
		} else if i2 == offIdx {
			i2++
		} else if i2 < offIdx {
			copy(seq[i1-delta:], seq[i2:offIdx])
			delta += (i2 - i1)
			i1, i2 = offIdx, offIdx+1
		} else {
			panic("wrong indices order")
		}
	}
	return delta, i1, i2
}

// GoError converts cmodule's Error to Go error and returns it.
func (errRun *Error) GoError(modules []Module) error {
	err := modules[errRun.Index].CToGoError(errRun.ModuleErr, errRun.SystemErr, errRun.Info)
	if err == nil {
		var errStr string
		if errRun.ModuleErr < 1000000 {
			errStr = "out of memory"
		} else {
			errStr = "unknown"
		}
		errStr = errStr + " (" + strconv.FormatInt(errRun.ModuleErr, 10)
		if errRun.SystemErr == 0 {
			errStr = errStr + ")"
		} else {
			errStr = errStr + ", " + strconv.FormatInt(errRun.SystemErr, 10) + ")"
		}
		if len(errRun.Info) > 0 {
			errStr = errStr + "; " + errRun.Info
		}
		err = errors.New(errStr)
	}
	return err
}

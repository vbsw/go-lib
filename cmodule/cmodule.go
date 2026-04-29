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
	sequenceChunks = 4
	MaxSequenceLen = min(uint64(C.SIZE_MAX)/sequenceChunks, uint64(maxInt/sequenceChunks))
)

// SequenceType represents processing type.
type SequenceType int

// Error is returned by function Process.
type Error struct {
	ModuleErr int64
	SystemErr int64
	Info      string
	Index     int
}

// Sequence is a buffer allocated in C.
// It holds pointers to functions and data.
type Sequence []C.uintptr_t

// Module provides abstraction to C functions and data.
type Module interface {
	// CProcessor returns pointer to a C function and C data.
	CProcessor(seqenceType SequenceType) (C.uintptr_t, C.uintptr_t)
	CToError(moduleErr, systemErr int64, info string) error
	SetCData(data C.uintptr_t)
}

// Len returns the number of modules in the Sequence.
func (seq Sequence) Len() int {
	return len(seq) / sequenceChunks
}

// NewSequence returns a new instance of Sequence.
func NewSequence(length int) Sequence {
	if length > 0 && uint64(length) <= MaxSequenceLen {
		var dataC *C.uintptr_t
		lengthTotal := length * sequenceChunks
		C.cmodule_alloc(&dataC, C.size_t(lengthTotal))
		if dataC != nil {
			return unsafe.Slice(dataC, lengthTotal)
		}
		return nil
	}
	panic("sequence length unsupported")
}

// Set sets functions and data for Process. Affects all if no indices given.
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

// Disable sets functions to be skipped in Process. Affects all if no indices given.
func (seq Sequence) Disable(indices ...int) {
	length := seq.Len()
	if len(indices) == 0 {
		for i := 0; i < length; i++ {
			seq[i], seq[i+length] = 0, 0
		}
	} else {
		for _, index := range indices {
			seq[index], seq[index+length] = 0, 0
		}
	}
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
		panic("passes number unsupported")
	}
	return nil
}

// Release releases C memory. Returns always nil.
func (seq Sequence) Release() Sequence {
	if len(seq) > 0 {
		C.cmodule_free(&seq[0])
	}
	return nil
}

func (errRun *Error) ToError(modules []Module) error {
	err := modules[errRun.Index].CToError(errRun.ModuleErr, errRun.SystemErr, errRun.Info)
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

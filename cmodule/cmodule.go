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
	maxInt   = int((^uint(0)) >> 1)
)

var (
	ErrBufferOutOfMemory   = errors.New("C buffer out of memory")
	ErrMaxModulesExceeded  = errors.New("maximum modules count exceeded")
	ErrDataModulesMismatch = errors.New("mismatch of data length and modules count")
)

type SequenceType int

type Sequence struct {
	Err        error
	modules    []Module
	buffer     []unsafe.Pointer
	syncLen    int
	bufferSize int
}

type Module interface {
	CProcessor(seqenceType SequenceType) (unsafe.Pointer, unsafe.Pointer)
	CToError(modErrId, sysErrId int64, errInfo string) error
	SetCData(data unsafe.Pointer)
}

// Len returns current number of modules in Sequence.
func (seq *Sequence) Len() int {
	return len(seq.modules)
}

// Cap returns possible number of modules in Sequence before reallocating.
func (seq *Sequence) Cap() int {
	return len(seq.buffer) / 2
}

// BufferSize returns C buffer size in bytes that's allocated by Sequence.
func (seq *Sequence) BufferSize() int {
	return seq.bufferSize
}

// EnsureCap preallocates memory where capacity is the number of modules
// that can be used before reallocating.
// Error is stored in Err.
func (seq *Sequence) EnsureCap(capacity int) error {
	if seq.Err == nil {
		if seq.Cap() < capacity {
			if int64(maxInt) >= (int64(capacity)+2)*2 {
				seq.initBuffer(capacity)
			} else {
				seq.Err = ErrMaxModulesExceeded
			}
		}
	}
	return seq.Err
}

// Add adds a module to Sequence.
// Error is stored in Err.
func (seq *Sequence) Add(module Module) error {
	if seq.Err == nil {
		if int64(maxInt) >= (int64(len(seq.modules))+2)*2 {
			seq.modules = append(seq.modules, module)
		} else {
			seq.Err = ErrMaxModulesExceeded
		}
	}
	return seq.Err
}

// Set prepares C functions and C data to be ran.
// Error is stored in Err.
func (seq *Sequence) Set(seqenceType SequenceType) error {
	if seq.Err == nil {
		if seq.Cap() < len(seq.modules) {
			seq.initBuffer(len(seq.modules))
		}
		if seq.Err == nil {
			for i, mod := range seq.modules {
				seq.buffer[i], seq.buffer[i+len(seq.modules)] = mod.CProcessor(seqenceType)
			}
		}
	}
	return seq.Err
}

// Run processes data in C.
// Error is stored in Err.
func (seq *Sequence) Run(passes int) error {
	if seq.Err == nil {
		var errIdx C.int32_t
		var err1C, err2C C.int64_t
		var errStrC *C.char
		lenC, sizeC, passesC := C.int32_t(len(seq.buffer)), C.int32_t(seq.bufferSize), C.int32_t(passes)
		C.vbsw_cmodule_proc(&seq.buffer[0], lenC, &sizeC, passesC, &errIdx, &err1C, &err2C, &errStrC)
		seq.bufferSize = int(sizeC)
		if err1C != 0 {
			seq.Err = seq.toError(int(errIdx), int64(err1C), int64(err2C), errStrC)
		}
	}
	return seq.Err
}

// SyncData writes C data to modules.
// Error is stored in Err.
func (seq *Sequence) SyncData() error {
	if seq.Err == nil {
		if seq.syncLen == len(seq.modules) {
			for i, mod := range seq.modules {
				mod.SetCData(seq.buffer[i+len(seq.modules)])
			}
		} else {
			seq.Err = ErrDataModulesMismatch
		}
	}
	return seq.Err
}

func (seq *Sequence) initBuffer(modulesLen int) {
	var bufferC *unsafe.Pointer
	lenC, sizeC := C.int32_t(len(seq.buffer)), C.int32_t(seq.bufferSize)
	if lenC > 0 {
		bufferC = &seq.buffer[0]
	}
	C.vbsw_cmodule_alloc_buffer(&bufferC, &lenC, &sizeC, C.int32_t(modulesLen))
	if sizeC > 0 {
		seq.buffer = unsafe.Slice(bufferC, int(lenC))
		seq.bufferSize = int(sizeC)
		seq.syncLen = modulesLen
	} else {
		seq.Err = ErrBufferOutOfMemory
	}
}

// Close releases modules and C memory.
func (seq *Sequence) Close() error {
	if seq.Err == nil {
		if len(seq.buffer) > 0 {
			C.vbsw_cmodule_free(&seq.buffer[0], C.int32_t(len(seq.buffer)))
			seq.modules = nil
			seq.buffer = nil
			seq.syncLen = 0
			seq.bufferSize = 0
		}
	}
	return seq.Err
}

func (seq *Sequence) toError(errIdx int, modErrId, sysErrId int64, errStrC *C.char) error {
	var errInfo string
	if errStrC != nil {
		errInfo = C.GoString(errStrC)
	}
	err := seq.modules[errIdx].CToError(modErrId, sysErrId, errInfo)
	if err == nil {
		var errStr string
		if modErrId < 1000000 {
			errStr = "out of memory"
		} else {
			errStr = "unknown"
		}
		errStr = errStr + " (" + strconv.FormatInt(modErrId, 10)
		if sysErrId == 0 {
			errStr = errStr + ")"
		} else {
			errStr = errStr + ", " + strconv.FormatInt(sysErrId, 10) + ")"
		}
		if len(errInfo) > 0 {
			errStr = errStr + "; " + errInfo
		}
		err = errors.New(errStr)
	}
	return err
}

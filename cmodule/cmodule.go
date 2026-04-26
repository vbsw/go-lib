/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package cmodule provides C batch processing.
package cmodule

// #include <stdlib.h>
// #include "cmodule.h"
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

var (
	ErrSeqBufferOutOfMemory = errors.New("sequence buffer out of memory")
)

const (
	maxInt64 = int64((^uint64(0)) >> 1)
	maxInt   = int((^uint(0)) >> 1)
)

type Module interface {
	CData() unsafe.Pointer
	CDestroyFunc() unsafe.Pointer
	CRunFunc() unsafe.Pointer
	CInitFunc() unsafe.Pointer
	CRemoveFunc() unsafe.Pointer
	CToError(modErrId, sysErrId int64, errInfo string) error
	SetCData(data unsafe.Pointer)
}

type Sequence struct {
	Modules []Module
	buffer  []unsafe.Pointer
}

// BufferSize returns size of buffer being used by Sequence.
func (seq *Sequence) BufferSize() int {
	return len(seq.buffer)
}

// ProcessInit processes Modules using CInitFunc of each Module.
// If no error occurs SetCData is called on each Module afterwards.
func (seq *Sequence) ProcessInit(passes int) error {
	if len(seq.Modules) > 0 {
		err := seq.ensureBuffer()
		if err == nil {
			var errIdx C.size_t
			var err1C, err2C C.longlong
			var errStrC *C.char
			modulesLenC, passesC := C.size_t(len(seq.Modules)), C.int(passes)
			for i, mod := range seq.Modules {
				seq.buffer[i] = mod.CInitFunc()
				seq.buffer[i+len(seq.Modules)] = mod.CData()
			}
			C.vbsw_cmodule_proc(&seq.buffer[0], modulesLenC, passesC, &errIdx, &err1C, &err2C, &errStrC)
			if err1C == 0 {
				for i, mod := range seq.Modules {
					mod.SetCData(seq.buffer[i+len(seq.Modules)])
				}
			} else {
				err = seq.toError(int(errIdx), int64(err1C), int64(err2C), errStrC)
			}
		}
		return err
	}
	return nil
}

// ProcessRun processes Modules using CRunFunc of each Module.
func (seq *Sequence) ProcessRun(passes int) error {
	if len(seq.Modules) > 0 {
		err := seq.ensureBuffer()
		if err == nil {
			var errIdx C.size_t
			var err1C, err2C C.longlong
			var errStrC *C.char
			modulesLenC, passesC := C.size_t(len(seq.Modules)), C.int(passes)
			for i, mod := range seq.Modules {
				seq.buffer[i] = mod.CRunFunc()
				seq.buffer[i+len(seq.Modules)] = mod.CData()
			}
			C.vbsw_cmodule_proc(&seq.buffer[0], modulesLenC, passesC, &errIdx, &err1C, &err2C, &errStrC)
			if err1C != 0 {
				err = seq.toError(int(errIdx), int64(err1C), int64(err2C), errStrC)
			}
		}
		return err
	}
	return nil
}

// ProcessRemove processes Modules using CRemoveFunc of each Module.
func (seq *Sequence) ProcessRemove(passes int, moduleIndices ...int) error {
	if len(seq.Modules) > 0 {
		err := seq.ensureBuffer()
		if err == nil {
			var errIdx C.size_t
			var err1C, err2C C.longlong
			var errStrC *C.char
			modulesLenC, passesC := C.size_t(len(seq.Modules)), C.int(passes)
			for i, mod := range seq.Modules {
				seq.buffer[i] = mod.CRemoveFunc()
				seq.buffer[i+len(seq.Modules)] = mod.CData()
			}
			C.vbsw_cmodule_rm(&seq.buffer[0], modulesLenC, passesC, &errIdx, &err1C, &err2C, &errStrC)
			if err1C != 0 {
				err = seq.toError(int(errIdx), int64(err1C), int64(err2C), errStrC)
			}
		}
		return err
	}
	return nil
}

// ProcessDestroy processes Modules using CDestroyFunc of each Module.
func (seq *Sequence) ProcessDestroy(passes int) error {
	if len(seq.Modules) > 0 {
		err := seq.ensureBuffer()
		if err == nil {
			var errIdx C.size_t
			var err1C, err2C C.longlong
			var errStrC *C.char
			modulesLenC, passesC := C.size_t(len(seq.Modules)), C.int(passes)
			for i, mod := range seq.Modules {
				seq.buffer[i] = mod.CDestroyFunc()
				seq.buffer[i+len(seq.Modules)] = mod.CData()
			}
			C.vbsw_cmodule_proc(&seq.buffer[0], modulesLenC, passesC, &errIdx, &err1C, &err2C, &errStrC)
			if err1C != 0 {
				err = seq.toError(int(errIdx), int64(err1C), int64(err2C), errStrC)
			}
		}
		return err
	}
	return nil
}

// Close releases buffer memory used by Sequence.
func (seq *Sequence) Close() error {
	if len(seq.buffer) > 0 {
		C.free(unsafe.Pointer(&seq.buffer[0]))
		seq.buffer = seq.buffer[:0]
	}
	return nil
}

func (seq *Sequence) ensureBuffer() error {
	if len(seq.buffer) < len(seq.Modules)*2 {
		var sizeC C.size_t
		var bufferC *unsafe.Pointer
		if len(seq.buffer) == 0 {
			C.vbsw_cmodule_alloc_buffer(&bufferC, &sizeC, nil, C.size_t(len(seq.Modules)))
		} else {
			C.vbsw_cmodule_alloc_buffer(&bufferC, &sizeC, &seq.buffer[0], C.size_t(len(seq.Modules)))
		}
		if sizeC > 0 && uint64(sizeC) <= uint64(maxInt) {
			seq.buffer = unsafe.Slice(bufferC, int(sizeC))
			return nil
		}
		return ErrSeqBufferOutOfMemory
	}
	return nil
}

func (seq *Sequence) toError(errIdx int, modErrId, sysErrId int64, errStrC *C.char) error {
	var errInfo string
	if errStrC != nil {
		errInfo = C.GoString(errStrC)
	}
	err := seq.Modules[errIdx].CToError(modErrId, sysErrId, errInfo)
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

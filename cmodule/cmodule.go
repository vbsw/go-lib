/*
 *          Copyright 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package cmodule provides C batch processing.
package cmodule

// #include "cmodule.h"
import "C"
import (
	"unsafe"
)

const (
	maxInt64 = int64((^uint64(0)) >> 1)
	maxInt32 = int32((^uint32(0)) >> 1)
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
	Modules    []Module
	buffer     unsafe.Pointer
	bufferSize int
}

// BufferSize returns size of buffer being used by the Sequence.
func (seq *Sequence) BufferSize() int {
	return seq.bufferSize
}

// ProcessInit processes Modules using CInitFunc of each Module.
// If no error occurs SetCData is called on each Module afterwards.
func (seq *Sequence) ProcessInit(passes int) error {
	return nil
}

// ProcessRun processes Modules using CRunFunc of each Module.
func (seq *Sequence) ProcessRun(passes int) error {
	return nil
}

// ProcessRemove processes Modules using CRemoveFunc of each Module.
func (seq *Sequence) ProcessRemove(passes int, moduleIndices ...int) error {
	return nil
}

// ProcessDestroy processes Modules using CDestroyFunc of each Module.
func (seq *Sequence) ProcessDestroy(passes int) error {
	return nil
}

// Close calls ProcessDestroy and releases buffer used by the Sequence
// if no error occurs.
func (seq *Sequence) Close() error {
	return nil
}

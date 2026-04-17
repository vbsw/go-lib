/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

//go:build windows
// +build windows

package fs

import (
	"os"
	"syscall"
)

// IsHidden returns true, if file is hidden.
func IsHidden(path string) bool {
	ptr, errPtr := syscall.UTF16PtrFromString(path)
	if errPtr == nil && ptr != nil {
		attrs, errAttrs := syscall.GetFileAttributes(ptr)
		if errAttrs == nil && attrs != syscall.INVALID_FILE_ATTRIBUTES {
			return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
		}
	}
	return false
}

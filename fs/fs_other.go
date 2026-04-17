/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

//go:build !windows
// +build !windows

package fs

import (
	"path/filepath"
	"strings"
)

// IsHidden returns true, if file is hidden.
func IsHidden(path string) bool {
	fileName := filepath.Base(path)
	return strings.HasPrefix(fileName, ".")
}

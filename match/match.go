/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package match provides simple wildcard string matching.
package match

import (
	"bytes"
	"unsafe"
)

type Operator int

type stateType int

const (
	And Operator = iota
	Or
	Xor
)

const (
	none stateType = iota
	skipping
	skippingEscape
	escape
)

// WildcardMatch returns true, if pattern matches str.
// Valid wildcards are "*" (any characters) and "?" (any single character).
// Escape character is "\".
func WildcardMatch(pattern, str string) bool {
	if len(str) > 0 {
		if len(pattern) > 0 {
			i, j, state := 0, 0, none
			for i < len(pattern) && j < len(str) {
				p, s := pattern[i], str[j]
				switch state {
				case none:
					if p == '*' {
						i, state = i+1, skipping
					} else if p == '?' {
						i, j = i+1, j+1
					} else if p == '\\' {
						i, state = i+1, escape
					} else if p != s {
						return false
					} else {
						i, j = i+1, j+1
					}
				case skipping:
					if p == '*' {
						i++
					} else if p == '\\' {
						i, state = i+1, skippingEscape
					} else if p == '?' {
						i, j = i+1, j+1
					} else if p != s {
						j++
					} else {
						i, j, state = i+1, j+1, none
					}
				case skippingEscape:
					if p == s {
						i, j, state = i+1, j+1, none
					} else {
						j++
					}
				case escape:
					if p == s {
						i, j, state = i+1, j+1, none
					} else {
						return false
					}
				}
			}
			if i == len(pattern) {
				return j == len(str) || state == skipping || state == skippingEscape
			} else {
				if j == len(str) {
					return pattern[i] == '\\' && i == len(pattern)
				} else {
					return false
				}
			}
		} else {
			return false
		}
	} else if len(pattern) > 0 {
		for _, b := range pattern {
			if b != '*' {
				return false
			}
		}
	}
	return true
}

// Contains returns true, if subslices are present in data.
//
// The evaluation is controlled by the provided operator:
//   - And: returns true if all subslices are contained in data.
//   - Or:  returns true if at least one of the subslices is contained in data.
//   - Xor: returns true if exclusivily one of the subslices is contained in data.
func Contains[D ~string | ~[]byte, S ~[]string | ~[][]byte](data D, subslices S, op Operator) bool {
	switch dataImpl := any(data).(type) {
	case []byte:
		switch subslicesImpl := any(subslices).(type) {
		case [][]byte:
			return containsBytes(dataImpl, subslicesImpl, op)
		case []string:
			return containsStrings(dataImpl, subslicesImpl, op)
		}
	case string:
		switch subslicesImpl := any(subslices).(type) {
		case [][]byte:
			return containsBytes(*(*[]byte)(unsafe.Pointer(&dataImpl)), subslicesImpl, op)
		case []string:
			return containsStrings(*(*[]byte)(unsafe.Pointer(&dataImpl)), subslicesImpl, op)
		}
	}
	return false
}

func containsBytes(data []byte, subslices [][]byte, op Operator) bool {
	switch op {
	case And:
		for _, slice := range subslices {
			if !bytes.Contains(data, slice) {
				return false
			}
		}
		return true
	case Or:
		for _, slice := range subslices {
			if bytes.Contains(data, slice) {
				return true
			}
		}
		return false
	case Xor:
		for i, slice := range subslices {
			if bytes.Contains(data, slice) {
				return i+1 == len(subslices) || !containsBytes(data, subslices[i+1:], Or)
			}
		}
		return false
	}
	return false
}

func containsStrings(data []byte, subslices []string, op Operator) bool {
	switch op {
	case And:
		for _, slice := range subslices {
			if !bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&slice))) {
				return false
			}
		}
		return true
	case Or:
		for _, slice := range subslices {
			if bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&slice))) {
				return true
			}
		}
		return false
	case Xor:
		for i, slice := range subslices {
			if bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&slice))) {
				return i+1 == len(subslices) || !containsStrings(data, subslices[i+1:], Or)
			}
		}
		return false
	}
	return false
}

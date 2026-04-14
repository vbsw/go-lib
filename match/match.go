/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package match provides simple wildcard string matching and substring check.
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

// WildcardMatch returns whether pattern matches string s.
// Possible wildcards are "*" (any characters) and "?" (any single character).
// Escape character is "\".
func WildcardMatch(pattern, s string) bool {
	if len(s) > 0 {
		if len(pattern) > 0 {
			i, j, state := 0, 0, none
			for i < len(pattern) && j < len(s) {
				pByte, sByte := pattern[i], s[j]
				switch state {
				case none:
					if pByte == '*' {
						i, state = i+1, skipping
					} else if pByte == '?' {
						i, j = i+1, j+1
					} else if pByte == '\\' {
						i, state = i+1, escape
					} else if pByte != sByte {
						return false
					} else {
						i, j = i+1, j+1
					}
				case skipping:
					if pByte == '*' {
						i++
					} else if pByte == '\\' {
						i, state = i+1, skippingEscape
					} else if pByte == '?' {
						i, j = i+1, j+1
					} else if pByte != sByte {
						j++
					} else {
						i, j, state = i+1, j+1, none
					}
				case skippingEscape:
					if pByte == sByte {
						i, j, state = i+1, j+1, none
					} else {
						j++
					}
				case escape:
					if pByte == sByte {
						i, j, state = i+1, j+1, none
					} else {
						return false
					}
				}
			}
			if i == len(pattern) {
				return j == len(s) || state == skipping || state == skippingEscape
			} else {
				if j == len(s) {
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

// Contains returns true, if substrings are present in string s.
//
// The evaluation is controlled by the provided operator:
//   - And: returns true if all substrings are contained in string s.
//   - Or:  returns true if at least one of the substrings is contained in string s.
//   - Xor: returns true if exclusivily one of the substrings is contained in string s.
func Contains[D ~string | ~[]byte, S ~[]string | ~[][]byte](s D, substrings S, op Operator) bool {
	switch data := any(s).(type) {
	case []byte:
		switch slices := any(substrings).(type) {
		case [][]byte:
			return containsBytes(data, slices, op)
		case []string:
			return containsStrings(data, slices, op)
		}
	case string:
		switch slices := any(substrings).(type) {
		case [][]byte:
			return containsBytes(*(*[]byte)(unsafe.Pointer(&data)), slices, op)
		case []string:
			return containsStrings(*(*[]byte)(unsafe.Pointer(&data)), slices, op)
		}
	}
	return false
}

func containsBytes(data []byte, slices [][]byte, op Operator) bool {
	switch op {
	case And:
		for _, slice := range slices {
			if !bytes.Contains(data, slice) {
				return false
			}
		}
		return true
	case Or:
		for _, slice := range slices {
			if bytes.Contains(data, slice) {
				return true
			}
		}
		return false
	case Xor:
		for i, slice := range slices {
			if bytes.Contains(data, slice) {
				return i+1 == len(slices) || !containsBytes(data, slices[i+1:], Or)
			}
		}
		return false
	}
	return false
}

func containsStrings(data []byte, slices []string, op Operator) bool {
	switch op {
	case And:
		for _, slice := range slices {
			if !bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&slice))) {
				return false
			}
		}
		return true
	case Or:
		for _, slice := range slices {
			if bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&slice))) {
				return true
			}
		}
		return false
	case Xor:
		for i, slice := range slices {
			if bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&slice))) {
				return i+1 == len(slices) || !containsStrings(data, slices[i+1:], Or)
			}
		}
		return false
	}
	return false
}

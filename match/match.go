/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package match provides simple wildcard string matching.
package match

type stateType int

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

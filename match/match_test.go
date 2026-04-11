/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package match

import (
	"path/filepath"
	"testing"
)

func TestWildcardMatch1(t *testing.T) {
	if !WildcardMatch("", "") {
		t.Error("failed pattern \"\" with \" \"")
	}
	if !WildcardMatch(" ", " ") {
		t.Error("failed pattern \" \" with \" \"")
	}
	if WildcardMatch("", "abcd") {
		t.Error("failed pattern \"\" with \"abcd\"")
	}
	if !WildcardMatch("*", "abcd") {
		t.Error("failed pattern \"*\" with \"abcd\"")
	}
	if !WildcardMatch("*", "") {
		t.Error("failed pattern \"*\" with \"\"")
	}
	if !WildcardMatch("****", "") {
		t.Error("failed pattern \"****\" with \"\"")
	}
	if WildcardMatch("*?", "") {
		t.Error("failed pattern \"*?\" with \"\"")
	}
	if !WildcardMatch("?*?", "abcd") {
		t.Error("failed pattern \"?*?\" with \"abcd\"")
	}
	if !WildcardMatch("?*?*?", "abc") {
		t.Error("failed pattern \"?*?*?\" with \"abc\"")
	}
}

func TestWildcardMatch2(t *testing.T) {
	if !WildcardMatch("a*d", "abcd") {
		t.Error("failed pattern \"a*d\" with \"abcd\"")
	}
	if WildcardMatch("???", "abcd") {
		t.Error("failed pattern \"???\" with \"abcd\"")
	}
	if !WildcardMatch("????", "abcd") {
		t.Error("failed pattern \"????\" with \"abcd\"")
	}
	if WildcardMatch("?????", "abcd") {
		t.Error("failed pattern \"?????\" with \"abcd\"")
	}
	if !WildcardMatch("?b*", "abcd") {
		t.Error("failed pattern \"?b*\" with \"abcd\"")
	}
	if !WildcardMatch("*c?", "abcd") {
		t.Error("failed pattern \"*c?\" with \"abcd\"")
	}
	if !WildcardMatch("a?c?e", "abcde") {
		t.Error("failed pattern \"a?c?e\" with \"abcde\"")
	}
	if !WildcardMatch("a?*e?*", "abcdefg") {
		t.Error("failed pattern \"a?*e?*\" with \"abcdefg\"")
	}
	if WildcardMatch(" a?*e?*", "abcdefg") {
		t.Error("failed pattern \" a?*e?*\" with \"abcdefg\"")
	}
}

func TestWildcardMatch3(t *testing.T) {
	if !WildcardMatch("*\\*d", "abc*d") {
		t.Error("failed pattern \"*\\*d\" with \"abc*d\"")
	}
	if !WildcardMatch("\\*\\?\\\\", "*?\\") {
		t.Error("failed pattern \"\\*\\?\\\\\" with \"*?\\\"")
	}
	if !WildcardMatch("*\\?x", "abcd?x") {
		t.Error("failed pattern \"*\\?x\" with \"abcd?x\"")
	}
	if !WildcardMatch("*\\?*", "abcd?efgh") {
		t.Error("failed pattern \"*\\?*\" with \"abcd?efgh\"")
	}
	if !WildcardMatch("abc\\\\", "abc\\") {
		t.Error("failed pattern \"abc\\\\\" with \"abc\\\"")
	}
	if WildcardMatch("abc\\", "abc\\") {
		t.Error("failed pattern \"abc\\\" with \"abc\\\"")
	}
	if !WildcardMatch("*\\", "abc\\") {
		t.Error("failed pattern \"*\\\" with \"abc\\\"")
	}
	if WildcardMatch("abc\\", "abc") {
		t.Error("failed pattern \"abc\\\" with \"abc\"")
	}
}

func BenchmarkWildcardMatch(b *testing.B) {
	result, str := true, "abcdefghijklmnopqrstuvwxyz"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = result && (WildcardMatch("*", str) == true)
		result = result && (WildcardMatch("*asdf", str) == false)
		result = result && (WildcardMatch("asdf*", str) == false)
		result = result && (WildcardMatch("abcdefghijklmnopqr*", str) == true)
		result = result && (WildcardMatch("*jklmnopqrstuvwxyz", str) == true)
		result = result && (WildcardMatch("*fghijklmnopqrst*", str) == true)
		result = result && (WildcardMatch("*fghijklmnopqrst", str) == false)
		result = result && (WildcardMatch("*efghijklm?opqrstuvwxyz", str) == true)
	}
	b.StopTimer()
	if !result {
		b.Fatal("wrong result")
	}
}

func BenchmarkFilepathMatch(b *testing.B) {
	result, str := true, "abcdefghijklmnopqrstuvwxyz"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		match, err := filepath.Match("*", str)
		result = result && err == nil && match == true
		match, err = filepath.Match("*asdf", str)
		result = result && err == nil && match == false
		match, err = filepath.Match("asdf*", str)
		result = result && err == nil && match == false
		match, err = filepath.Match("abcdefghijklmnopqr*", str)
		result = result && err == nil && match == true
		match, err = filepath.Match("*jklmnopqrstuvwxyz", str)
		result = result && err == nil && match == true
		match, err = filepath.Match("*fghijklmnopqrst*", str)
		result = result && err == nil && match == true
		match, err = filepath.Match("*fghijklmnopqrst", str)
		result = result && err == nil && match == false
		match, err = filepath.Match("*efghijklm?opqrstuvwxyz", str)
		result = result && err == nil && match == true
	}
	b.StopTimer()
	if !result {
		b.Fatal("wrong result")
	}
}

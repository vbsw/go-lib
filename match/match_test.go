/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package match

import (
	"bytes"
	"path/filepath"
	"testing"
	"unsafe"
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

func TestContainsAnd(t *testing.T) {
	data := []byte("abcdefghijklmnopqrstuvwxyhallozabcdefghijklmiddlenopqrstuvwxyzabciaodefghijklmnopqrstuvwxyzend")
	slices := []string{"hallo", "middle", "ciao", "end", "abcdefgh-none", ""}
	if !Contains(data, slices[:1], And) {
		t.Error("failed \"hallo\"")
	}
	if !Contains(data, slices[1:2], And) {
		t.Error("failed \"middle\"")
	}
	if !Contains(data, slices[2:3], And) {
		t.Error("failed \"ciao\"")
	}
	if !Contains(data, slices[:4], And) {
		t.Error("failed slices[:4]")
	}
	if !Contains(data, slices[5:6], And) {
		t.Error("failed \"\"")
	}
	if Contains(data, slices[4:5], And) {
		t.Error("failed \"abcdefgh-none\"")
	}
	if Contains(data, slices, And) {
		t.Error("failed slices")
	}
}

func TestContainsXor(t *testing.T) {
	data := []byte("abcdefghijklmnopqrstuvwxyhallozabcdefghijklmiddlenopqrstuvwxyzabciaodefghijklmnopqrstuvwxyzend")
	slices := []string{"hallo", "middle", "ciaoXX", "end", "abcdefgh-none", ""}
	if !Contains(data, slices[3:5], Xor) {
		t.Error("failed \"slices[3:5]\"")
	}
	if !Contains(data, slices[2:5], Xor) {
		t.Error("failed \"slices[2:5]\"")
	}
	if Contains(data, slices, Xor) {
		t.Error("failed slices")
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

func BenchmarkContains(b *testing.B) {
	result := true
	data := "abcdefghijklmnopqrstuvwxyhallozabcdefghijklmiddlenopqrstuvwxyzabciaodefghijklmnopqrstuvwxyzend"
	slices := []string{"middle", "hallo", "end", "ciao", "abcdefgh-none"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = result && (Contains(data, slices, And) == false)
	}
	b.StopTimer()
	if !result {
		b.Fatal("wrong result")
	}
}

func BenchmarkContainsStd(b *testing.B) {
	result := true
	data := []byte("abcdefghijklmnopqrstuvwxyhallozabcdefghijklmiddlenopqrstuvwxyzabciaodefghijklmnopqrstuvwxyzend")
	slices := []string{"middle", "hallo", "end", "ciao", "abcdefgh-none"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultTmp := true
		for _, str := range slices {
			if !bytes.Contains(data, *(*[]byte)(unsafe.Pointer(&str))) {
				resultTmp = false
				break
			}
		}
		result = result && (resultTmp == false)
	}
	b.StopTimer()
	if !result {
		b.Fatal("wrong result")
	}
}

/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package cl

import "testing"

func TestSingleArgA(t *testing.T) {
	cmdLine := New([]string{"asdf", "--version"})
	versionArg := cmdLine.Search("-v", "--version")
	if !versionArg.Available() {
		t.Error()
	} else {
		if versionArg.Keys[0] != cmdLine.Arguments[1] {
			t.Error(versionArg.Keys[0])
		}
	}
	if cmdLine.Matched[0] {
		t.Error()
	}
	if !cmdLine.Matched[1] {
		t.Error()
	}
}

func TestSingleArgB(t *testing.T) {
	cmdLine := New([]string{"--start", "asdf", "-s", "qwer"})
	startArg := cmdLine.Search("-s", "--start")

	if !startArg.Available() {
		t.Error()
	} else {
		if startArg.Count() == 2 {
			if startArg.Keys[0] != cmdLine.Arguments[0] {
				t.Error(startArg.Keys[0])
			}
			if startArg.Keys[1] != cmdLine.Arguments[2] {
				t.Error(startArg.Keys[1])
			}
		} else {
			t.Error(startArg.Count())
		}
	}
}

func TestArgByDelimiterA(t *testing.T) {
	cmdLine := New([]string{"asdf", "--start=123"})
	cmdLine.Delimiter = NewDelimiter("=", "")
	startArg := cmdLine.SearchByDelimiter("-s", "--start")

	if !startArg.Available() {
		t.Error()
	} else {
		if startArg.Count() == 1 {
			if startArg.Values[0] != "123" {
				t.Error(startArg.Values[0])
			}
		} else {
			t.Error(startArg.Count())
		}
	}
	if cmdLine.Matched[0] {
		t.Error()
	}
	if !cmdLine.Matched[1] {
		t.Error()
	}
}

func TestArgByDelimiterB(t *testing.T) {
	cmdLine := New([]string{"asdf", "--start", "123"})
	cmdLine.Delimiter = NewDelimiter("=", " ")
	startArg := cmdLine.SearchByDelimiter("-s", "--start")

	if !startArg.Available() {
		t.Error()
	} else {
		if startArg.Count() == 1 {
			if startArg.Values[0] != "123" {
				t.Error(startArg.Values[0])
			}
		} else {
			t.Error(startArg.Count())
		}
	}
	if cmdLine.Matched[0] {
		t.Error()
	}
	if !cmdLine.Matched[1] {
		t.Error()
	}
}

func TestArgByDelimiterC(t *testing.T) {
	cmdLine := New([]string{"asdf", "--start123"})
	cmdLine.Delimiter = NewDelimiter("=", "")
	startArg := cmdLine.SearchByDelimiter("-s", "--start")

	if !startArg.Available() {
		t.Error()
	} else {
		if startArg.Count() == 1 {
			if startArg.Values[0] != "123" {
				t.Error(startArg.Values[0])
			}
		} else {
			t.Error(startArg.Count())
		}
	}
	if cmdLine.Matched[0] {
		t.Error()
	}
	if !cmdLine.Matched[1] {
		t.Error()
	}
}

func TestRest(t *testing.T) {
	cmdLine := New([]string{"--start", "asdf", "-s", "qwer"})
	cmdLine.Search("--start", "-s")
	unmatched := cmdLine.Unmatched()

	if !cmdLine.Matched[0] {
		t.Error(cmdLine.Matched[0])

	} else if !cmdLine.Matched[2] {
		t.Error(cmdLine.Matched[0])

	} else if cmdLine.Matched[1] || cmdLine.Matched[3] || cmdLine.Matched[4] {
		t.Error("wrongly matched")

	} else if len(unmatched.Keys) != 2 {
		t.Error(len(unmatched.Keys))

	} else if unmatched.Keys[0] != cmdLine.Arguments[1] {
		t.Error(unmatched.Keys[0])

	} else if unmatched.Keys[1] != cmdLine.Arguments[3] {
		t.Error(unmatched.Keys[1])
	}
	cmdLine.Search("asdf", "qwer")
	if !cmdLine.Matched[4] || cmdLine.Unmatched() != nil {
		t.Error("wrongly not matched")
	}
}

func TestSpace(t *testing.T) {
	delimiter := NewDelimiter(" ", "asdf", " ")
	if len(delimiter.Separators) != 1 {
		t.Error(len(delimiter.Separators))
	} else if delimiter.Separators[0] != "asdf" {
		t.Error(delimiter.Separators[0])
	} else if !delimiter.HasSpaceSeparator {
		t.Error()
	}
	delimiter = NewDelimiter(" ", "asdf")
	if len(delimiter.Separators) != 1 {
		t.Error(len(delimiter.Separators))
	} else if delimiter.Separators[0] != "asdf" {
		t.Error(delimiter.Separators[0])
	} else if !delimiter.HasSpaceSeparator {
		t.Error()
	}
	delimiter = NewDelimiter("asdf", " ")
	if len(delimiter.Separators) != 1 {
		t.Error(len(delimiter.Separators))
	} else if delimiter.Separators[0] != "asdf" {
		t.Error(delimiter.Separators[0])
	} else if !delimiter.HasSpaceSeparator {
		t.Error()
	}
	delimiter = NewDelimiter("asdf", " ", "qwer")
	if len(delimiter.Separators) != 2 {
		t.Error(len(delimiter.Separators))
	} else if delimiter.Separators[0] != "asdf" || delimiter.Separators[1] != "qwer" {
		t.Error(delimiter.Separators[0], delimiter.Separators[1])
	} else if !delimiter.HasSpaceSeparator {
		t.Error()
	}
}

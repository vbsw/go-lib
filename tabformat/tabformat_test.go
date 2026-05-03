/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package tabformat

import (
	"testing"
)

func TestLineBeginEnd1(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\nccc ddd\reee fff\r\nggg\r")
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 1 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 0 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 7 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 8 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestLineBeginEnd2(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\nccc ddd\reee fff\r\nggg\r")
	parser.Next(line)
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 2 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 8 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 15 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 16 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestLineBeginEnd3(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\nccc ddd\reee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 3 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 16 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 23 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 25 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestLineBeginEnd4(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\nccc ddd\reee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	parser.Next(line)
	success := parser.Next(line)
	if success != false {
		t.Error("wrong success:", success)
	} else {
		parser.IgnoreOpenEnd = true
		success = parser.Next(line)
		if success != true {
			t.Error("wrong success:", success)
		} else if parser.LineNumber != 4 {
			t.Error("line number wrong:", parser.LineNumber)
		} else if parser.LineBegin >= parser.LineEnd {
			t.Error("empty line:", parser.LineBegin, parser.LineEnd)
		} else if parser.LineBegin != 25 {
			t.Error("wrong line begin:", parser.LineBegin)
		} else if parser.LineEnd != 28 {
			t.Error("wrong line end:", parser.LineEnd)
		} else if parser.nextLineBegin != 29 {
			t.Error("wrong next line begin:", parser.nextLineBegin)
		}
	}
}

func TestElementBeginEnd(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\nccc ddd\reee fff\r\nggg\r")
	parser.Next(line)
	if parser.KeyBegin != 0 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.KeyEnd != 3 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if parser.Indent != 0 {
		t.Error("wrong indent:", parser.Indent)
	} else if string(line[parser.KeyBegin:parser.KeyEnd]) != "aaa" {
		t.Error("wrong key:", string(line[parser.KeyBegin:parser.KeyEnd]))
	} else if string(line[parser.ValBegin:parser.ValEnd]) != "bbb" {
		t.Error("wrong value:", string(line[parser.ValBegin:parser.ValEnd]))
	} else {
		parser.Next(line)
		if string(line[parser.KeyBegin:parser.KeyEnd]) != "ccc" {
			t.Error("wrong key:", string(line[parser.KeyBegin:parser.KeyEnd]))
		} else if string(line[parser.ValBegin:parser.ValEnd]) != "ddd" {
			t.Error("wrong value:", string(line[parser.ValBegin:parser.ValEnd]))
		}
	}
}

func TestIndent(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\n\tccc ddd\r\t\teee fff\r\n\t\t\t\tggg\r")
	parser.Next(line) // aaa bbb
	if parser.Indent != 0 {
		t.Error("wrong indent:", parser.Indent)
	} else {
		parser.Next(line) // ccc ddd
		if parser.Indent != 1 {
			t.Error("wrong indent:", parser.Indent)
		} else {
			parser.Next(line) // eee fff
			if parser.Indent != 2 {
				t.Error("wrong indent:", parser.Indent)
			} else {
				parser.IgnoreOpenEnd = true
				parser.Next(line) // tggg
				if parser.Indent != 4 {
					t.Error("wrong indent:", parser.Indent)
				}
			}
		}
	}
}

func TestElement(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa   bbb   |   ccc  ddd eee   \r ")
	parser.Next(line)
	if len(parser.Value(line)) != 3 {
		t.Error("wrong value length:", len(parser.Value(line)))
	} else if string(parser.Value(line)) != "bbb" {
		t.Error("wrong value:", string(parser.Value(line)))
	} else {
		parser.Next(line)
		if len(parser.Key(line)) != 3 {
			t.Error("wrong key length:", len(parser.Key(line)))
		} else if string(parser.Key(line)) != "ccc" {
			t.Error("wrong key:", string(parser.Value(line)))
		} else if len(parser.Value(line)) != 7 {
			t.Error("wrong value length:", len(parser.Value(line)))
		} else if string(parser.Value(line)) != "ddd eee" {
			t.Error("wrong value:", string(parser.Value(line)))
		}
	}
}

func TestComment1(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb#ccc ddd|eee fff\r\ngg#g\r")
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineBegin != 0 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 23 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 25 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	} else if parser.KeyBegin != 0 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.KeyEnd != 3 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if parser.ValBegin != 4 {
		t.Error("wrong value begin:", parser.ValBegin)
	} else if parser.ValEnd != 7 {
		t.Error("wrong value end:", parser.ValEnd)
	}
}

func TestComment2(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb#ccc ddd|eee fff\r\ngg#g\r")
	parser.Next(line)
	parser.IgnoreOpenEnd = true
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineBegin != 25 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 29 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 30 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	} else if parser.KeyBegin != 25 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.KeyEnd != 27 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if parser.ValBegin != 27 {
		t.Error("wrong value begin:", parser.ValBegin)
	} else if parser.ValEnd != 27 {
		t.Error("wrong value end:", parser.ValEnd)
	}
}

func TestComment3(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb#ccc ddd|eee fff\r\ngg#g\r")
	parser.Next(line)
	if parser.state != stateNewLine {
		t.Error("wrong state:", parser.state)
	}
	parser.IgnoreOpenEnd = true
	parser.Next(line)
	if parser.state != stateNewLine {
		t.Error("wrong state:", parser.state)
	}
	success := parser.Next(line)
	if success != false {
		t.Error("wrong success:", success)
	} else if parser.LineBegin != 25 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 29 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 30 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	} else if parser.KeyBegin != 25 || parser.KeyEnd != 27 {
		t.Error("key wrong:", parser.KeyBegin, parser.KeyEnd)
	} else if parser.ValBegin != 27 || parser.ValEnd != 27 {
		t.Error("value wrong:", parser.ValBegin, parser.ValEnd)
	}
}

func TestInline1(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 1 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.Indent != 0 {
		t.Error("wrong indent:", parser.Indent)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 0 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 23 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 25 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestInline2(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	parser.Next(line)
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 1 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.Indent != 1 {
		t.Error("wrong indent:", parser.Indent)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 0 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 23 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 25 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestInline3(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 1 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.Indent != 1 {
		t.Error("wrong indent:", parser.Indent)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 0 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 23 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 25 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestInline4(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	parser.Next(line)
	parser.IgnoreOpenEnd = true
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else if parser.LineNumber != 2 {
		t.Error("line number wrong:", parser.LineNumber)
	} else if parser.Indent != 0 {
		t.Error("wrong indent:", parser.Indent)
	} else if parser.LineBegin >= parser.LineEnd {
		t.Error("empty line:", parser.LineBegin, parser.LineEnd)
	} else if parser.LineBegin != 25 {
		t.Error("wrong line begin:", parser.LineBegin)
	} else if parser.LineEnd != 28 {
		t.Error("wrong line end:", parser.LineEnd)
	} else if parser.nextLineBegin != 29 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	}
}

func TestInlineElement1(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	if parser.KeyBegin != 8 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.KeyEnd != 11 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if parser.ValBegin != 12 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.ValEnd != 15 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if string(line[parser.KeyBegin:parser.KeyEnd]) != "ccc" {
		t.Error("wrong key:", string(line[parser.KeyBegin:parser.KeyEnd]))
	} else if string(line[parser.ValBegin:parser.ValEnd]) != "ddd" {
		t.Error("wrong value:", string(line[parser.ValBegin:parser.ValEnd]))
	}
}

func TestInlineElement2(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	parser.Next(line)
	if parser.KeyBegin != 16 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.KeyEnd != 19 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if parser.ValBegin != 20 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.ValEnd != 23 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if string(line[parser.KeyBegin:parser.KeyEnd]) != "eee" {
		t.Error("wrong key:", string(line[parser.KeyBegin:parser.KeyEnd]))
	} else if string(line[parser.ValBegin:parser.ValEnd]) != "fff" {
		t.Error("wrong value:", string(line[parser.ValBegin:parser.ValEnd]))
	}
}

func TestInlineElement3(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb\\ccc ddd|eee fff\r\nggg\r")
	parser.Next(line)
	parser.Next(line)
	parser.Next(line)
	parser.IgnoreOpenEnd = true
	parser.Next(line)
	if parser.KeyBegin != 25 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.KeyEnd != 28 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if parser.ValBegin != 28 {
		t.Error("wrong key begin:", parser.KeyBegin)
	} else if parser.ValEnd != 28 {
		t.Error("wrong key end:", parser.KeyEnd)
	} else if string(line[parser.KeyBegin:parser.KeyEnd]) != "ggg" {
		t.Error("wrong key:", string(line[parser.KeyBegin:parser.KeyEnd]))
	} else if string(line[parser.ValBegin:parser.ValEnd]) != "" {
		t.Error("wrong value:", string(line[parser.ValBegin:parser.ValEnd]))
	}
}

func TestIncomplete1(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb")
	success := parser.Next(line)
	if success != false {
		t.Error("wrong success:", success)
	} else if parser.nextLineBegin != 0 {
		t.Error("wrong next line begin:", parser.nextLineBegin)
	} else if parser.Rest(len(line)) != 7 {
		t.Error("wrong incomplete length:", parser.Rest(len(line)))
	} else {
		parser.IgnoreOpenEnd = true
		success = parser.Next(line)
		if success != true {
			t.Error("wrong success:", success)
		} else if parser.nextLineBegin != 7 {
			t.Error("wrong next line begin:", parser.nextLineBegin)
		} else if parser.Rest(len(line)) != 0 {
			t.Error("wrong incomplete length:", parser.Rest(len(line)))
		}
	}
}

func TestIncomplete2(t *testing.T) {
	var parser ByteParser
	lineAAA := []byte("aaa bbb")
	lineCCC := []byte("aaa bbb\nccc ddd")
	parser.Next(lineAAA)
	success := parser.Next(lineCCC)
	if success != true {
		t.Error("wrong success:", success)
	} else {
		success = parser.Next(lineCCC)
		if success != false {
			t.Error("wrong success:", success)
		} else {
			rest := parser.Reset(len(lineCCC))
			parser.IgnoreOpenEnd = true
			lineEEE := []byte("aaa bbb\nccc ddd|eee fff\r")
			lineEEE = lineEEE[len(lineCCC)-rest:]
			success = parser.Next(lineEEE) // stateInlineChild
			if success != true {
				t.Error("wrong success:", success)
			} else if string(lineEEE[parser.KeyBegin:parser.KeyEnd]) != "ccc" {
				t.Error("wrong key:", string(lineEEE[parser.KeyBegin:parser.KeyEnd]))
			} else if string(lineEEE[parser.ValBegin:parser.ValEnd]) != "ddd" {
				t.Error("wrong key:", string(lineEEE[parser.ValBegin:parser.ValEnd]))
			}
		}
	}
}

func TestIncomplete3(t *testing.T) {
	var parser ByteParser
	line := []byte("aaa bbb#ccc ddd|eee fff\r\ngg#g\r")
	success := parser.Next(line)
	if success != true {
		t.Error("wrong success:", success)
	} else {
		success = parser.Next(line)
		if success != false {
			t.Error("wrong success:", success)
		} else if parser.Rest(len(line)) != 5 {
			t.Error("wrong incomplete length:", parser.Rest(len(line)))
		} else {
			parser.Reset(len(line))
			success = parser.Next(line)
			if success != true {
				t.Error("wrong success:", success)
			}
		}
	}
}

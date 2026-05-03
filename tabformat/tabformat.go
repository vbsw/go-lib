/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package tabformat provides a parser for a simple, tab-indented data format.
// The format represents hierarchical structures using leading tab characters.
package tabformat

type stateType uint8

const (
	stateNewLine stateType = iota
	stateNewLinePrefix
	stateInlineChild
	stateInlineSibling
	stateComment
)

// ByteParser holds parse state.
type ByteParser struct {
	KeyBegin, KeyEnd   int
	ValBegin, ValEnd   int
	Indent, LineNumber int
	LineBegin, LineEnd int
	nextKeyBegin       int
	nextLineBegin      int
	state              stateType
	IgnoreOpenEnd      bool
}

// Next reads bytes and stores key and value.
// Returns true if line has been read.
func (p *ByteParser) Next(bytes []byte) bool {
	for true {
		switch p.state {
		case stateNewLine:
			if p.parseLineBounds(bytes) {
				p.parseIndentation(bytes)
				p.state = stateNewLinePrefix
			} else {
				return false
			}
		case stateNewLinePrefix:
			p.KeyBegin = iSkipWhitespaceAndCharB(p.KeyBegin, p.LineEnd, bytes, '|')
			if isCommentB(bytes[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				if p.parseInlineChildPrefix(bytes) {
					p.state = stateNewLinePrefix
				} else {
					p.parseKeyValue(bytes)
					return true
				}
			}
		case stateInlineChild:
			p.KeyBegin = iSkipWhitespaceB(p.nextKeyBegin, p.LineEnd, bytes)
			for p.parseInlineChildPrefix(bytes) {
				p.KeyBegin = iSkipWhitespaceB(p.nextKeyBegin, p.LineEnd, bytes)
			}
			if isCommentB(bytes[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				p.Indent++
				p.parseKeyValue(bytes)
				return true
			}
		case stateInlineSibling:
			p.KeyBegin = iSkipWhitespaceAndCharB(p.ValEnd, p.LineEnd, bytes, '|')
			if isCommentB(bytes[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				p.parseKeyValue(bytes)
				return true
			}
		}
	}
	return false
}

// Reset sets all members except LineNumber to zero.
// Returns unparsed number of bytes.
func (p *ByteParser) Reset(total int) int {
	rest := total - p.nextLineBegin
	p.KeyBegin, p.KeyEnd = 0, 0
	p.ValBegin, p.ValEnd = 0, 0
	p.LineBegin, p.LineEnd = 0, 0
	p.nextKeyBegin, p.nextLineBegin = 0, 0
	p.state, p.IgnoreOpenEnd = stateNewLine, false
	p.Indent = 0
	return rest
}

// Rest returns unparsed number of bytes.
func (p *ByteParser) Rest(total int) int {
	return total - p.nextLineBegin
}

// Key returns key slice.
func (p *ByteParser) Key(bytes []byte) []byte {
	return bytes[p.KeyBegin:p.KeyEnd]
}

// Value returns value slice.
func (p *ByteParser) Value(bytes []byte) []byte {
	return bytes[p.ValBegin:p.ValEnd]
}

func (p *ByteParser) parseLineBounds(bytes []byte) bool {
	for i := p.nextLineBegin; i < len(bytes); i++ {
		if bytes[i] == '\r' {
			if i1 := i + 1; i1 < len(bytes) {
				p.LineNumber++
				p.LineBegin, p.LineEnd = p.nextLineBegin, i
				if bytes[i1] == '\n' {
					p.nextLineBegin = i1 + 1
				} else {
					p.nextLineBegin = i1
				}
				return true
			} else if p.IgnoreOpenEnd {
				p.LineNumber++
				p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, i, i+1
				return true
			}
			return false
		} else if bytes[i] == '\n' {
			p.LineNumber++
			p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, i, i+1
			return true
		}
	}
	if p.nextLineBegin < len(bytes) && p.IgnoreOpenEnd {
		p.LineNumber++
		p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, len(bytes), len(bytes)
		return true
	}
	return false
}

func (p *ByteParser) parseIndentation(bytes []byte) {
	p.KeyBegin, p.Indent = p.LineBegin, 0
	for p.KeyBegin < p.LineEnd && bytes[p.KeyBegin] == '\t' {
		p.Indent++
		p.KeyBegin++
	}
}

func (p *ByteParser) parseInlineChildPrefix(bytes []byte) bool {
	if p.KeyBegin < p.LineEnd && bytes[p.KeyBegin] == '\\' {
		if keyBegin1 := p.KeyBegin + 1; keyBegin1 < p.LineEnd {
			if iByte := bytes[keyBegin1]; iByte != '\\' && iByte != '#' && iByte != '|' {
				p.KeyBegin, p.Indent = p.KeyBegin+2, p.Indent+1
				return true
			}
		}
	}
	return false
}

func (p *ByteParser) parseKeyValue(bytes []byte) {
	stateOld := p.state
	p.KeyEnd, p.state = iParseKey(p.KeyBegin, p.LineEnd, bytes, p.state)
	if stateOld == p.state {
		p.ValBegin = iSkipWhitespaceB(p.KeyEnd, p.LineEnd, bytes)
		p.ValEnd, p.state = iParseValue(p.ValBegin, p.LineEnd, bytes, p.state)
		p.nextKeyBegin = p.ValEnd + 1
		p.ValEnd = iSkipWhitespaceReverseB(p.ValBegin, p.ValEnd, bytes)
		if stateOld == p.state {
			p.state = stateNewLine
		}
	} else {
		p.ValBegin = p.KeyEnd
		p.ValEnd = p.KeyEnd
		p.nextKeyBegin = p.KeyEnd + 1
	}
	p.KeyEnd = iSkipWhitespaceReverseB(p.KeyBegin, p.KeyEnd, bytes)
}

func iParseKey(begin, end int, bytes []byte, state stateType) (int, stateType) {
	var escape bool
	for i := begin; i < end; i++ {
		if iByte := bytes[i]; iByte < 0 || iByte > 32 { // non whitespace
			if iByte == '\\' {
				escape = !escape
			} else if iByte == '#' {
				if escape {
					escape = false
				} else {
					return i, stateNewLine
				}
			} else if iByte == '|' {
				if escape {
					escape = false
				} else {
					return i, stateInlineSibling
				}
			} else if escape {
				return i - 1, stateInlineChild
			}
		} else { // whitespace
			return i, state
		}
	}
	return end, state
}

func iParseValue(begin, end int, bytes []byte, state stateType) (int, stateType) {
	var escape bool
	for i := begin; i < end; i++ {
		if iByte := bytes[i]; iByte < 0 || iByte > 32 { // non whitespace
			if iByte == '\\' {
				escape = !escape
			} else if iByte == '#' {
				if escape {
					escape = false
				} else {
					return i, stateNewLine
				}
			} else if iByte == '|' {
				if escape {
					escape = false
				} else {
					return i, stateInlineSibling
				}
			} else if escape {
				return i - 1, stateInlineChild
			}
		} else if escape {
			return i - 1, stateInlineChild
		}
	}
	return end, state
}

func (p *ByteParser) iInlineChild(bytes []byte) stateType {
	for i := p.ValEnd + 1; i < len(bytes); i++ {
		if iByte := bytes[i]; iByte > 0 && iByte < 33 || iByte == '|' {
			i++
		}
	}
	return 0
}

func isCommentB(bytes []byte) bool {
	for iByte := range bytes {
		if iByte < 0 || iByte > 32 { // non whitespace
			if iByte == '#' {
				return true
			}
			return false
		} // else: whitespace
	}
	return false
}

func iSkipWhitespaceB(i, lineEnd int, bytes []byte) int {
	for i < lineEnd {
		if iByte := bytes[i]; iByte > 0 && iByte < 33 {
			i++
		} else {
			break
		}
	}
	return i
}

func iSkipWhitespaceReverseB(left, right int, bytes []byte) int {
	i := right - 1
	for i >= left {
		if iByte := bytes[i]; iByte > 0 && iByte < 33 {
			i--
		} else {
			break
		}
	}
	return i + 1
}

func iSkipWhitespaceAndCharB(i, lineEnd int, bytes []byte, charToSkip byte) int {
	for i < lineEnd {
		if iByte := bytes[i]; iByte > 0 && iByte < 33 || iByte == charToSkip {
			i++
		} else {
			break
		}
	}
	return i
}

func iSkipWhitespaceAndFirstCharB(i, lineEnd int, bytes []byte, charToSkip byte) int {
	var charSkipped bool
	for i < lineEnd {
		iByte := bytes[i]
		if iByte > 0 && iByte < 33 {
			i++
		} else if iByte == charToSkip {
			if !charSkipped {
				charSkipped = true
				i++
			} else {
				break
			}
		} else {
			break
		}
	}
	return i
}

func isCommentS(str string) bool {
	for i := 0; i < len(str); i++ {
		iByte := str[i]
		if iByte < 0 || iByte > 32 {
			if iByte == '#' {
				return true
			}
			return false
		} // else: skip whitespace
	}
	return false
}

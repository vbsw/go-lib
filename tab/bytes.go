/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package tabformat

// ByteParser holds parse state.
type ByteParser struct {
	KeyBegin, KeyEnd   int
	ValBegin, ValEnd   int
	KeyLen, ValLen     int
	Indent, LineNumber int
	LineBegin, LineEnd int
	LineLen            int
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
				p.LineLen = p.LineEnd - p.LineBegin
				p.parseIndentation(bytes)
				p.state = stateNewLinePrefix
			} else {
				return false
			}
		case stateNewLinePrefix:
			p.KeyBegin = iSkipWhitespaceAndCharB(bytes, p.KeyBegin, p.LineEnd, '|')
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
			p.KeyBegin = iSkipWhitespaceB(bytes, p.nextKeyBegin, p.LineEnd)
			for p.parseInlineChildPrefix(bytes) {
				p.KeyBegin = iSkipWhitespaceB(bytes, p.KeyBegin, p.LineEnd)
			}
			if isCommentB(bytes[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				p.Indent++
				p.parseKeyValue(bytes)
				return true
			}
		case stateInlineSibling:
			p.KeyBegin = iSkipWhitespaceAndCharB(bytes, p.ValEnd, p.LineEnd, '|')
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
	*p = ByteParser{LineNumber: p.LineNumber}
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

// Line returns line slice.
func (p *ByteParser) Line(bytes []byte) []byte {
	return bytes[p.LineBegin:p.LineEnd]
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
	p.KeyEnd, p.state = iParseKeyB(bytes, p.KeyBegin, p.LineEnd, p.state)
	if stateOld == p.state {
		p.ValBegin = iSkipWhitespaceB(bytes, p.KeyEnd, p.LineEnd)
		p.ValEnd, p.state = iParseValueB(bytes, p.ValBegin, p.LineEnd, p.state)
		p.nextKeyBegin = p.ValEnd + 1
		p.ValEnd = iSkipWhitespaceReverseB(bytes, p.ValBegin, p.ValEnd)
		if stateOld == p.state {
			p.state = stateNewLine
		}
	} else {
		p.ValBegin = p.KeyEnd
		p.ValEnd = p.KeyEnd
		p.nextKeyBegin = p.KeyEnd + 1
	}
	p.KeyEnd = iSkipWhitespaceReverseB(bytes, p.KeyBegin, p.KeyEnd)
	p.KeyLen = p.KeyEnd - p.KeyBegin
	p.ValLen = p.ValEnd - p.ValBegin
}

func iParseKeyB(bytes []byte, from, to int, state stateType) (int, stateType) {
	var escape bool
	for i := from; i < to; i++ {
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
	return to, state
}

func iParseValueB(bytes []byte, from, to int, state stateType) (int, stateType) {
	var escape bool
	for i := from; i < to; i++ {
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
	return to, state
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

func iSkipWhitespaceB(bytes []byte, from, to int) int {
	for i := from; i < to; i++ {
		if iByte := bytes[i]; iByte < 0 || iByte > 32 {
			return i
		}
	}
	return to
}

func iSkipWhitespaceReverseB(bytes []byte, from, to int) int {
	for i := to - 1; i > from; i-- {
		if iByte := bytes[i]; iByte < 0 || iByte > 32 {
			return i + 1
		}
	}
	if iByte := bytes[from]; iByte < 0 || iByte > 32 {
		return from + 1
	}
	return from
}

func iSkipWhitespaceAndCharB(bytes []byte, from, to int, charToSkip byte) int {
	for i := from; i < to; i++ {
		if iByte := bytes[i]; (iByte < 0 || iByte > 32) && iByte != charToSkip {
			return i
		}
	}
	return to
}

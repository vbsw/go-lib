/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package tabformat

// StringParser holds parse state.
type StringParser struct {
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
func (p *StringParser) Next(s string) bool {
	for true {
		switch p.state {
		case stateNewLine:
			if p.parseLineBounds(s) {
				p.LineLen = p.LineEnd - p.LineBegin
				p.parseIndentation(s)
				p.state = stateNewLinePrefix
			} else {
				return false
			}
		case stateNewLinePrefix:
			p.KeyBegin = iSkipWhitespaceAndCharS(s, p.KeyBegin, p.LineEnd, '|')
			if isCommentS(s[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				if p.parseInlineChildPrefix(s) {
					p.state = stateNewLinePrefix
				} else {
					p.parseKeyValue(s)
					return true
				}
			}
		case stateInlineChild:
			p.KeyBegin = iSkipWhitespaceS(s, p.nextKeyBegin, p.LineEnd)
			for p.parseInlineChildPrefix(s) {
				p.KeyBegin = iSkipWhitespaceS(s, p.KeyBegin, p.LineEnd)
			}
			if isCommentS(s[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				p.Indent++
				p.parseKeyValue(s)
				return true
			}
		case stateInlineSibling:
			p.KeyBegin = iSkipWhitespaceAndCharS(s, p.ValEnd, p.LineEnd, '|')
			if isCommentS(s[p.KeyBegin:p.LineEnd]) {
				p.state = stateNewLine
			} else {
				p.parseKeyValue(s)
				return true
			}
		}
	}
	return false
}

// Reset sets all members except LineNumber to zero.
// Returns unparsed number of bytes.
func (p *StringParser) Reset(total int) int {
	rest := total - p.nextLineBegin
	*p = StringParser{LineNumber: p.LineNumber}
	return rest
}

// Rest returns unparsed number of bytes.
func (p *StringParser) Rest(total int) int {
	return total - p.nextLineBegin
}

// Key returns key string.
func (p *StringParser) Key(s string) string {
	return s[p.KeyBegin:p.KeyEnd]
}

// Value returns value string.
func (p *StringParser) Value(s string) string {
	return s[p.ValBegin:p.ValEnd]
}

// Line returns line string.
func (p *StringParser) Line(s string) string {
	return s[p.LineBegin:p.LineEnd]
}

func (p *StringParser) parseLineBounds(s string) bool {
	for i := p.nextLineBegin; i < len(s); i++ {
		if s[i] == '\r' {
			if i1 := i + 1; i1 < len(s) {
				p.LineNumber++
				p.LineBegin, p.LineEnd = p.nextLineBegin, i
				if s[i1] == '\n' {
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
		} else if s[i] == '\n' {
			p.LineNumber++
			p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, i, i+1
			return true
		}
	}
	if p.nextLineBegin < len(s) && p.IgnoreOpenEnd {
		p.LineNumber++
		p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, len(s), len(s)
		return true
	}
	return false
}

func (p *StringParser) parseIndentation(s string) {
	p.KeyBegin, p.Indent = p.LineBegin, 0
	for p.KeyBegin < p.LineEnd && s[p.KeyBegin] == '\t' {
		p.Indent++
		p.KeyBegin++
	}
}

func (p *StringParser) parseInlineChildPrefix(s string) bool {
	if p.KeyBegin < p.LineEnd && s[p.KeyBegin] == '\\' {
		if keyBegin1 := p.KeyBegin + 1; keyBegin1 < p.LineEnd {
			if iByte := s[keyBegin1]; iByte != '\\' && iByte != '#' && iByte != '|' {
				p.KeyBegin, p.Indent = p.KeyBegin+2, p.Indent+1
				return true
			}
		}
	}
	return false
}

func (p *StringParser) parseKeyValue(s string) {
	stateOld := p.state
	p.KeyEnd, p.state = iParseKeyS(s, p.KeyBegin, p.LineEnd, p.state)
	if stateOld == p.state {
		p.ValBegin = iSkipWhitespaceS(s, p.KeyEnd, p.LineEnd)
		p.ValEnd, p.state = iParseValueS(s, p.ValBegin, p.LineEnd, p.state)
		p.nextKeyBegin = p.ValEnd + 1
		p.ValEnd = iSkipWhitespaceReverseS(s, p.ValBegin, p.ValEnd)
		if stateOld == p.state {
			p.state = stateNewLine
		}
	} else {
		p.ValBegin = p.KeyEnd
		p.ValEnd = p.KeyEnd
		p.nextKeyBegin = p.KeyEnd + 1
	}
	p.KeyEnd = iSkipWhitespaceReverseS(s, p.KeyBegin, p.KeyEnd)
	p.KeyLen = p.KeyEnd - p.KeyBegin
	p.ValLen = p.ValEnd - p.ValBegin
}

func iParseKeyS(s string, from, to int, state stateType) (int, stateType) {
	var escape bool
	for i := from; i < to; i++ {
		if iByte := s[i]; iByte < 0 || iByte > 32 { // non whitespace
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

func iParseValueS(s string, from, to int, state stateType) (int, stateType) {
	var escape bool
	for i := from; i < to; i++ {
		if iByte := s[i]; iByte < 0 || iByte > 32 { // non whitespace
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

func isCommentS(str string) bool {
	for i := 0; i < len(str); i++ {
		if iByte := str[i]; iByte < 0 || iByte > 32 { // non whitespace
			if iByte == '#' {
				return true
			}
			return false
		} // else: whitespace
	}
	return false
}

func iSkipWhitespaceS(s string, from, to int) int {
	for i := from; i < to; i++ {
		if iByte := s[i]; iByte < 0 || iByte > 32 {
			return i
		}
	}
	return to
}

func iSkipWhitespaceReverseS(s string, from, to int) int {
	for i := to - 1; i > from; i-- {
		if iByte := s[i]; iByte < 0 || iByte > 32 {
			return i + 1
		}
	}
	if iByte := s[from]; iByte < 0 || iByte > 32 {
		return from + 1
	}
	return from
}

func iSkipWhitespaceAndCharS(s string, from, to int, charToSkip byte) int {
	for i := from; i < to; i++ {
		if iByte := s[i]; (iByte < 0 || iByte > 32) && iByte != charToSkip {
			return i
		}
	}
	return to
}

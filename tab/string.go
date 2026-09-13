/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package tab

// ParserS holds parse state. It provides parsing of string.
type ParserS struct {
	KeyBegin, KeyEnd   int
	ValBegin, ValEnd   int
	KeyLen, ValLen     int
	Indent, LineNumber int
	LineBegin, LineEnd int
	LineLen            int
	nextKeyBegin       int
	nextLineBegin      int
	state              stateType
	ParseLastLine      bool
}

// ListParserS holds parse state. It provides parsing of list values in a string.
type ListParserS struct {
	ListBegin, ListEnd    int
	EntryBegin, EntryEnd  int
	EntryLen              int
	SeparatorBegin, Index int
}

// Next reads bytes and stores key and value.
// Returns true if line has been read.
func (p *ParserS) Next(s string) bool {
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
func (p *ParserS) Reset(total int) int {
	rest := total - p.nextLineBegin
	*p = ParserS{LineNumber: p.LineNumber}
	return rest
}

// Rest returns unparsed number of bytes.
func (p *ParserS) Rest(total int) int {
	return total - p.nextLineBegin
}

// Key returns key string.
func (p *ParserS) Key(s string) string {
	return s[p.KeyBegin:p.KeyEnd]
}

// Value returns value string.
func (p *ParserS) Value(s string) string {
	return s[p.ValBegin:p.ValEnd]
}

// Line returns line string.
func (p *ParserS) Line(s string) string {
	return s[p.LineBegin:p.LineEnd]
}

// ListParserV returns list parser for key.
func (p *ParserS) ListParserK() ListParserS {
	var listParser ListParserS
	listParser.ListBegin = p.KeyBegin
	listParser.ListEnd = p.KeyEnd
	listParser.EntryBegin = p.KeyBegin
	listParser.EntryEnd = p.KeyBegin
	listParser.SeparatorBegin = p.KeyBegin - 1
	listParser.Index = -1
	return listParser
}

// ListParserV returns list parser for value.
func (p *ParserS) ListParserV() ListParserS {
	var listParser ListParserS
	listParser.ListBegin = p.ValBegin
	listParser.ListEnd = p.ValEnd
	listParser.EntryBegin = p.ValBegin
	listParser.EntryEnd = p.ValBegin
	listParser.SeparatorBegin = p.ValBegin - 1
	listParser.Index = -1
	return listParser
}

// Entry returns entry slice.
func (p *ListParserS) Entry(bytes []byte) []byte {
	return bytes[p.EntryBegin:p.EntryEnd]
}

// Init initializes list parser state.
func (p *ListParserS) Init(listBegin, listEnd int) {
	p.ListBegin = listBegin
	p.ListEnd = listEnd
	p.EntryBegin = listBegin
	p.EntryEnd = listBegin
	p.SeparatorBegin = listBegin - 1
	p.Index = -1
}

// Next reads bytes and stores entry offsets.
// Returns true if entry has been read.
func (p *ListParserS) Next(bytes []byte, separator byte) bool {
	for i := p.SeparatorBegin + 1; i < p.ListEnd; i++ {
		iByte := bytes[i]
		if iByte == separator && separator != ' ' {
			p.EntryBegin, p.EntryEnd, p.SeparatorBegin, p.EntryLen = i, i, i, 0
			p.Index++
			return true
		} else if iByte > 32 {
			p.EntryBegin = i
			p.SeparatorBegin = p.ListEnd
			for j := i + 1; j < p.ListEnd; j++ {
				if bytes[j] == separator {
					if separator == ' ' {
						nextEntryBegin := iSkipWhitespaceB(bytes, j+1, p.ListEnd)
						if nextEntryBegin < p.ListEnd {
							p.SeparatorBegin = nextEntryBegin - 1
						}
					} else {
						p.SeparatorBegin = j
					}
					break
				}
			}
			p.EntryEnd = iSkipWhitespaceReverseB(bytes, i+1, p.SeparatorBegin)
			p.EntryLen = p.EntryEnd - p.EntryBegin
			p.Index++
			return true
		}
	}
	if p.SeparatorBegin < p.ListEnd {
		p.EntryBegin, p.EntryEnd, p.SeparatorBegin, p.EntryLen = p.ListEnd, p.ListEnd, p.ListEnd, 0
		p.Index++
		return true
	}
	p.EntryBegin, p.EntryEnd, p.EntryLen = p.ListEnd, p.ListEnd, 0
	return false
}

func (p *ParserS) parseLineBounds(s string) bool {
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
			} else if p.ParseLastLine {
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
	if p.nextLineBegin < len(s) && p.ParseLastLine {
		p.LineNumber++
		p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, len(s), len(s)
		return true
	}
	return false
}

func (p *ParserS) parseIndentation(s string) {
	p.KeyBegin, p.Indent = p.LineBegin, 0
	for p.KeyBegin < p.LineEnd && s[p.KeyBegin] == '\t' {
		p.Indent++
		p.KeyBegin++
	}
}

func (p *ParserS) parseInlineChildPrefix(s string) bool {
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

func (p *ParserS) parseKeyValue(s string) {
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
		if iByte := s[i]; iByte > 32 {
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
		if iByte := s[i]; iByte > 32 {
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
		if iByte := str[i]; iByte > 32 {
			if iByte == '#' {
				return true
			}
			return false
		}
	}
	return false
}

func iSkipWhitespaceS(s string, from, to int) int {
	for i := from; i < to; i++ {
		if s[i] > 32 {
			return i
		}
	}
	return to
}

func iSkipWhitespaceReverseS(s string, from, to int) int {
	for i := to - 1; i >= from; i-- {
		if s[i] > 32 {
			return i + 1
		}
	}
	return from
}

func iSkipWhitespaceAndCharS(s string, from, to int, charToSkip byte) int {
	for i := from; i < to; i++ {
		if iByte := s[i]; iByte > 32 && iByte != charToSkip {
			return i
		}
	}
	return to
}

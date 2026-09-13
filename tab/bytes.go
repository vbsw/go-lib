/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package tab

// ParserB holds parse state. It provides parsing of byte array.
type ParserB struct {
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

// ListParserB holds parse state. It provides parsing of list values of byte array.
type ListParserB struct {
	ListBegin, ListEnd    int
	EntryBegin, EntryEnd  int
	EntryLen              int
	SeparatorBegin, Index int
}

// Next reads bytes and stores key and value offsets.
// Returns true if line has been read.
func (p *ParserB) Next(bytes []byte) bool {
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
func (p *ParserB) Reset(total int) int {
	rest := total - p.nextLineBegin
	*p = ParserB{LineNumber: p.LineNumber}
	return rest
}

// Rest returns unparsed number of bytes.
func (p *ParserB) Rest(total int) int {
	return total - p.nextLineBegin
}

// Key returns key slice.
func (p *ParserB) Key(bytes []byte) []byte {
	return bytes[p.KeyBegin:p.KeyEnd]
}

// Value returns value slice.
func (p *ParserB) Value(bytes []byte) []byte {
	return bytes[p.ValBegin:p.ValEnd]
}

// Line returns line slice.
func (p *ParserB) Line(bytes []byte) []byte {
	return bytes[p.LineBegin:p.LineEnd]
}

// ListParserV returns list parser for key.
func (p *ParserB) ListParserK() ListParserB {
	var listParser ListParserB
	listParser.ListBegin = p.KeyBegin
	listParser.ListEnd = p.KeyEnd
	listParser.EntryBegin = p.KeyBegin
	listParser.EntryEnd = p.KeyBegin
	listParser.SeparatorBegin = p.KeyBegin - 1
	listParser.Index = -1
	return listParser
}

// ListParserV returns list parser for value.
func (p *ParserB) ListParserV() ListParserB {
	var listParser ListParserB
	listParser.ListBegin = p.ValBegin
	listParser.ListEnd = p.ValEnd
	listParser.EntryBegin = p.ValBegin
	listParser.EntryEnd = p.ValBegin
	listParser.SeparatorBegin = p.ValBegin - 1
	listParser.Index = -1
	return listParser
}

// Entry returns entry slice.
func (p *ListParserB) Entry(bytes []byte) []byte {
	return bytes[p.EntryBegin:p.EntryEnd]
}

// Init initializes list parser state.
func (p *ListParserB) Init(listBegin, listEnd int) {
	p.ListBegin = listBegin
	p.ListEnd = listEnd
	p.EntryBegin = listBegin
	p.EntryEnd = listBegin
	p.SeparatorBegin = listBegin - 1
	p.Index = -1
}

// Next reads bytes and stores entry offsets.
// Returns true if entry has been read.
func (p *ListParserB) Next(bytes []byte, separator byte) bool {
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

func (p *ParserB) parseLineBounds(bytes []byte) bool {
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
			} else if p.ParseLastLine {
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
	if p.nextLineBegin < len(bytes) && p.ParseLastLine {
		p.LineNumber++
		p.LineBegin, p.LineEnd, p.nextLineBegin = p.nextLineBegin, len(bytes), len(bytes)
		return true
	}
	return false
}

func (p *ParserB) parseIndentation(bytes []byte) {
	p.KeyBegin, p.Indent = p.LineBegin, 0
	for p.KeyBegin < p.LineEnd && bytes[p.KeyBegin] == '\t' {
		p.Indent++
		p.KeyBegin++
	}
}

func (p *ParserB) parseInlineChildPrefix(bytes []byte) bool {
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

func (p *ParserB) parseKeyValue(bytes []byte) {
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
		if iByte := bytes[i]; iByte > 32 {
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
		if iByte := bytes[i]; iByte > 32 {
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
		if iByte > 32 {
			if iByte == '#' {
				return true
			}
			return false
		}
	}
	return false
}

func iSkipWhitespaceB(bytes []byte, from, to int) int {
	for i := from; i < to; i++ {
		if bytes[i] > 32 {
			return i
		}
	}
	return to
}

func iSkipWhitespaceReverseB(bytes []byte, from, to int) int {
	for i := to - 1; i >= from; i-- {
		if bytes[i] > 32 {
			return i + 1
		}
	}
	return from
}

func iSkipWhitespaceAndCharB(bytes []byte, from, to int, charToSkip byte) int {
	for i := from; i < to; i++ {
		if iByte := bytes[i]; iByte > 32 && iByte != charToSkip {
			return i
		}
	}
	return to
}

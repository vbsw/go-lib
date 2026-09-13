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

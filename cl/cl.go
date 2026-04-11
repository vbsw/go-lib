/*
 *       Copyright 2025, 2026, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package cl parses command line arguments.
package cl

import "strings"

// CommandLine represents command line.
type CommandLine struct {
	Arguments []string
	Matched   []bool
	Delimiter *Delimiter
}

// Arguments represents arguments returned by command line search.
type Arguments struct {
	Keys    []string
	Values  []string
	Indices []int
}

// Delimiter represents separators between key and value.
type Delimiter struct {
	Separators        []string
	HasSpaceSeparator bool
	HasEmptySeparator bool
}

// New returns a new instance of CommandLine.
func New(args []string, delimiter *Delimiter) *CommandLine {
	cmdLine := new(CommandLine)
	length := len(args)
	if length > 0 {
		cmdLine.Arguments = make([]string, length)
		copy(cmdLine.Arguments, args)
	}
	// last value means all arguments are matched
	cmdLine.Matched = make([]bool, length+1)
	cmdLine.Matched[length] = bool(length == 0)
	cmdLine.Delimiter = delimiter
	return cmdLine
}

// NewDelimiter returns a new instance of Delimiter. An empty separator ""
// sets the HasEmptySeparator flag, and the space separator " " sets
// the HasSpaceSeparator flag.
func NewDelimiter(separators ...string) *Delimiter {
	delimiter, firstEntry := new(Delimiter), true
	for _, separator := range separators {
		if separator == "" {
			delimiter.HasEmptySeparator = true
		} else if separator == " " {
			delimiter.HasSpaceSeparator = true
		} else {
			if firstEntry {
				delimiter.Separators = make([]string, 1)
				delimiter.Separators[0] = separator
				firstEntry = false
			} else {
				delimiter.Separators = append(delimiter.Separators, separator)
			}
		}
	}
	return delimiter
}

// Match searches flags in CommandLine.Arguments and returns them.
func (cmdLine *CommandLine) Match(flags ...string) *Arguments {
	var args *Arguments
	if cmdLine != nil {
		if !cmdLine.Matched[len(cmdLine.Arguments)] && len(flags) > 0 {
			allMatched := true
			for i, argument := range cmdLine.Arguments {
				if !cmdLine.Matched[i] {
					for _, searchedFlag := range flags {
						if argument == searchedFlag {
							if args == nil {
								args = new(Arguments)
							}
							args.Keys = append(args.Keys, searchedFlag)
							args.Indices = append(args.Indices, i)
							cmdLine.Matched[i] = true
							break
						}
					}
					allMatched = allMatched && cmdLine.Matched[i]
				}
			}
			cmdLine.Matched[len(cmdLine.Arguments)] = allMatched
		}
	}
	return args
}

// MatchDelimited searches flags in CommandLine.Arguments and returns them.
// The search considers a delimiter, that separates key and value within argument.
// Is Delimiter.HasSpaceSeparator set, then two arguments are treated as one argument
// with key and value separated by space.
func (cmdLine *CommandLine) MatchDelimited(flags ...string) *Arguments {
	var args *Arguments
	if cmdLine != nil {
		if !cmdLine.Matched[len(cmdLine.Arguments)] && len(flags) > 0 {
			if cmdLine.Delimiter.HasSpaceSeparator {
				args = cmdLine.matchDelimitedWithSpace(flags)
			} else {
				args = cmdLine.matchDelimitedWithoutSpace(flags)
			}
		}
	}
	return args
}

// RevertMatched sets arguments in CommandLine to unmatched using args.Indices.
func (cmdLine *CommandLine) RevertMatched(args ...*Arguments) {
	allMatched := cmdLine.Matched[len(cmdLine.Arguments)]
	for _, aargs := range args {
		if aargs != nil {
			for _, index := range aargs.Indices {
				cmdLine.Matched[index] = false
				allMatched = false
			}
		}
	}
	cmdLine.Matched[len(cmdLine.Arguments)] = allMatched
}

// Available returns true if at least one argument is available. If args == nil returns false.
func (args *Arguments) Available() bool {
	return args != nil && len(args.Keys) > 0
}

// Count returns number of arguments. If args == nil returns 0.
func (args *Arguments) Count() int {
	if args != nil {
		return len(args.Keys)
	}
	return 0
}

// HasKey returns true if key is in Keys. If args == nil returns false.
func (args *Arguments) HasKey(key string) bool {
	if args != nil {
		for _, k := range args.Keys {
			if k == key {
				return true
			}
		}
	}
	return false
}

// HasValue returns true if value is in Values. If args == nil returns false.
func (args *Arguments) HasValue(value string) bool {
	if args != nil {
		for _, v := range args.Values {
			if v == value {
				return true
			}
		}
	}
	return false
}

// HasIndex returns true if index is in Indices. If args == nil returns false.
func (args *Arguments) HasIndex(index int) bool {
	if args != nil {
		for _, i := range args.Indices {
			if i == index {
				return true
			}
		}
	}
	return false
}

// KeyAt returns key at index. If args == nil or index is not in list returns fallback.
func (args *Arguments) KeyAt(index int, fallback string) string {
	if args != nil {
		if index >= 0 && index < len(args.Keys) {
			return args.Keys[index]
		}
	}
	return fallback
}

// ValueAt returns value at index. If args == nil or index is not in list returns fallback.
func (args *Arguments) ValueAt(index int, fallback string) string {
	if args != nil {
		if index >= 0 && index < len(args.Values) {
			return args.Values[index]
		}
	}
	return fallback
}

// IndexAt returns index at index. If args == nil or index is not in list returns fallback.
func (args *Arguments) IndexAt(index int, fallback int) int {
	if args != nil {
		if index >= 0 && index < len(args.Indices) {
			return args.Indices[index]
		}
	}
	return fallback
}

// Unmatched returns arguments that haven't been matched.
func (cmdLine *CommandLine) Unmatched() *Arguments {
	var args *Arguments
	if cmdLine != nil {
		length := len(cmdLine.Arguments)
		if !cmdLine.Matched[length] {
			args = new(Arguments)
			for i, arg := range cmdLine.Arguments {
				if !cmdLine.Matched[i] {
					args.Keys = append(args.Keys, arg)
					args.Indices = append(args.Indices, i)
				}
			}
		}
	}
	return args
}

func (cmdLine *CommandLine) matchDelimitedWithSpace(flags []string) *Arguments {
	var args *Arguments
	length, allMatched := len(cmdLine.Arguments), true
	for i := 0; i < length; i++ {
		if !cmdLine.Matched[i] {
			argument := cmdLine.Arguments[i]
			for _, searchedFlag := range flags {
				if len(argument) > len(searchedFlag) {
					value, ok := cmdLine.Delimiter.argumentValue(argument, searchedFlag)
					if ok {
						if args == nil {
							args = new(Arguments)
						}
						args.Keys = append(args.Keys, searchedFlag)
						args.Values = append(args.Values, value)
						args.Indices = append(args.Indices, i)
						cmdLine.Matched[i] = true
						break
					}
				} else if argument == searchedFlag {
					iNxt := i + 1
					if args == nil {
						args = new(Arguments)
					}
					args.Keys = append(args.Keys, searchedFlag)
					args.Indices = append(args.Indices, i)
					cmdLine.Matched[i] = true
					if iNxt < length && !cmdLine.Matched[iNxt] {
						args.Values = append(args.Values, cmdLine.Arguments[iNxt])
						cmdLine.Matched[iNxt] = true
						i = iNxt
					} else {
						args.Values = append(args.Values, "")
					}
					break
				}
			}
			allMatched = allMatched && cmdLine.Matched[i]
		}
	}
	cmdLine.Matched[length] = allMatched
	return args
}

func (cmdLine *CommandLine) matchDelimitedWithoutSpace(flags []string) *Arguments {
	var args *Arguments
	allMatched := true
	for i, argument := range cmdLine.Arguments {
		if !cmdLine.Matched[i] {
			for _, searchedFlag := range flags {
				if len(argument) > len(searchedFlag) {
					value, ok := cmdLine.Delimiter.argumentValue(argument, searchedFlag)
					if ok {
						if args == nil {
							args = new(Arguments)
						}
						args.Keys = append(args.Keys, searchedFlag)
						args.Values = append(args.Values, value)
						args.Indices = append(args.Indices, i)
						cmdLine.Matched[i] = true
						break
					}
				} else if cmdLine.Delimiter.HasEmptySeparator && argument == searchedFlag {
					if args == nil {
						args = new(Arguments)
					}
					args.Keys = append(args.Keys, searchedFlag)
					args.Values = append(args.Values, "")
					args.Indices = append(args.Indices, i)
					cmdLine.Matched[i] = true
					break
				}
			}
			allMatched = allMatched && cmdLine.Matched[i]
		}
	}
	cmdLine.Matched[len(cmdLine.Arguments)] = allMatched
	return args
}

func (delimiter *Delimiter) argumentValue(argument, key string) (string, bool) {
	if strings.HasPrefix(argument, key) {
		argumentWithoutKey := argument[len(key):]
		for _, separator := range delimiter.Separators {
			if strings.HasPrefix(argumentWithoutKey, separator) {
				return argumentWithoutKey[len(separator):], true
			}
		}
		if delimiter.HasEmptySeparator {
			return argumentWithoutKey, true
		}
	}
	return "", false
}

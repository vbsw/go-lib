/*
 *          Copyright 2025, Vitali Baumtrok.
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
	Keys   []string
	Values []string
}

// Delimiter represents separators between key and value.
type Delimiter struct {
	Separators        []string
	HasSpaceSeparator bool
	HasEmptySeparator bool
}

// New returns a new instance of CommandLine.
func New(args []string) *CommandLine {
	cmdLine := new(CommandLine)
	length := len(args)
	if length > 0 {
		cmdLine.Arguments = make([]string, length)
		copy(cmdLine.Arguments, args)
	}
	// last value means all arguments are matched
	cmdLine.Matched = make([]bool, length+1)
	cmdLine.Matched[length] = bool(length == 0)
	return cmdLine
}

// NewDelimiter returns a new instance of Delimiter. An empty separator ""
// sets the HasEmptySeparator flag for delimiter, and the space separator " " sets
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

// Search compairs CommandLine.Arguments with searchTerms and returns matches.
func (cmdLine *CommandLine) Search(searchTerms ...string) *Arguments {
	var args *Arguments
	if cmdLine != nil {
		if !cmdLine.Matched[len(cmdLine.Arguments)] && len(searchTerms) > 0 {
			allMatched := true
			for i, argument := range cmdLine.Arguments {
				if !cmdLine.Matched[i] {
					for _, searchTerm := range searchTerms {
						if argument == searchTerm {
							if args == nil {
								args = new(Arguments)
							}
							args.Keys = append(args.Keys, searchTerm)
							cmdLine.Matched[i] = true
							break
						}
					}
					// check only last iteration
					allMatched = allMatched && cmdLine.Matched[i]
				}
			}
			cmdLine.Matched[len(cmdLine.Arguments)] = allMatched
		}
	}
	return args
}

// SearchByDelimiter compairs CommandLine.Arguments with searchTerms and returns matches.
// The search considers a delimiter, that separates key and value within parameter.
// Is Delimiter.HasSpaceSeparator set, then two arguments are treated as one argument
// with key and value separated by space.
func (cmdLine *CommandLine) SearchByDelimiter(searchTerms ...string) *Arguments {
	var args *Arguments
	if cmdLine != nil {
		if !cmdLine.Matched[len(cmdLine.Arguments)] && len(searchTerms) > 0 {
			if cmdLine.Delimiter.HasSpaceSeparator {
				args = cmdLine.searchPairsWithSpace(searchTerms)
			} else {
				args = cmdLine.searchPairsWithoutSpace(searchTerms)
			}
		}
	}
	return args
}

// Available returns true, if at least one argument is available.
func (args *Arguments) Available() bool {
	return args != nil && len(args.Keys) > 0
}

// Count returns number of arguments.
func (args *Arguments) Count() int {
	if args != nil {
		return len(args.Keys)
	}
	return 0
}

// Unmatched returns command line arguments that haven't been matched by the search.
func (cmdLine *CommandLine) Unmatched() *Arguments {
	var args *Arguments
	if cmdLine != nil {
		length := len(cmdLine.Arguments)
		if !cmdLine.Matched[length] {
			args = new(Arguments)
			for i, arg := range cmdLine.Arguments {
				if !cmdLine.Matched[i] {
					args.Keys = append(args.Keys, arg)
				}
			}
		}
	}
	return args
}

func (cmdLine *CommandLine) searchPairsWithSpace(searchTerms []string) *Arguments {
	var args *Arguments
	length, allMatched := len(cmdLine.Arguments), true
	for i := 0; i < length; i++ {
		if !cmdLine.Matched[i] {
			argument := cmdLine.Arguments[i]
			for _, searchTerm := range searchTerms {
				if len(argument) > len(searchTerm) {
					value, ok := cmdLine.Delimiter.argumentValue(argument, searchTerm)
					if ok {
						if args == nil {
							args = new(Arguments)
						}
						args.Keys = append(args.Keys, searchTerm)
						args.Values = append(args.Values, value)
						cmdLine.Matched[i] = true
						break
					}
				} else if argument == searchTerm {
					value, iNxt := "", i+1
					cmdLine.Matched[i] = true
					if iNxt < length && !cmdLine.Matched[iNxt] {
						value = cmdLine.Arguments[iNxt]
						cmdLine.Matched[iNxt] = true
						i = iNxt
					}
					if args == nil {
						args = new(Arguments)
					}
					args.Keys = append(args.Keys, searchTerm)
					args.Values = append(args.Values, value)
					break
				}
			}
			allMatched = allMatched && cmdLine.Matched[i]
		}
	}
	cmdLine.Matched[length] = allMatched
	return args
}

func (cmdLine *CommandLine) searchPairsWithoutSpace(searchTerms []string) *Arguments {
	var args *Arguments
	allMatched := true
	for i, argument := range cmdLine.Arguments {
		if !cmdLine.Matched[i] {
			for _, searchTerm := range searchTerms {
				if len(argument) > len(searchTerm) {
					value, ok := cmdLine.Delimiter.argumentValue(argument, searchTerm)
					if ok {
						if args == nil {
							args = new(Arguments)
						}
						args.Keys = append(args.Keys, searchTerm)
						args.Values = append(args.Values, value)
						cmdLine.Matched[i] = true
						break
					}
				} else if cmdLine.Delimiter.HasEmptySeparator && argument == searchTerm {
					if args == nil {
						args = new(Arguments)
					}
					args.Keys = append(args.Keys, searchTerm)
					args.Values = append(args.Values, "")
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

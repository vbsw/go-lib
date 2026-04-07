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
	Origins []int
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

// Search compairs CommandLine.Arguments with searchedKeys and returns matches.
func (cmdLine *CommandLine) Search(searchedKeys ...string) *Arguments {
	var args *Arguments
	if cmdLine != nil {
		if !cmdLine.Matched[len(cmdLine.Arguments)] && len(searchedKeys) > 0 {
			allMatched := true
			for i, argument := range cmdLine.Arguments {
				if !cmdLine.Matched[i] {
					for _, searchedKey := range searchedKeys {
						if argument == searchedKey {
							if args == nil {
								args = new(Arguments)
							}
							args.Keys = append(args.Keys, searchedKey)
							args.Origins = append(args.Origins, i)
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

// SearchByDelimiter compairs CommandLine.Arguments with searchedKeys and returns matches.
// The search considers a delimiter, that separates key and value within parameter.
// Is Delimiter.HasSpaceSeparator set, then two arguments are treated as one argument
// with key and value separated by space.
func (cmdLine *CommandLine) SearchByDelimiter(searchedKeys ...string) *Arguments {
	var args *Arguments
	if cmdLine != nil {
		if !cmdLine.Matched[len(cmdLine.Arguments)] && len(searchedKeys) > 0 {
			if cmdLine.Delimiter.HasSpaceSeparator {
				args = cmdLine.searchPairsWithSpace(searchedKeys)
			} else {
				args = cmdLine.searchPairsWithoutSpace(searchedKeys)
			}
		}
	}
	return args
}

// RevertMatched sets arguments in CommandLine to unmatched.
func (cmdLine *CommandLine) RevertMatched(argsList ...*Arguments) {
	for _, args := range argsList {
		if args != nil {
			for _, origin := range args.Origins {
				cmdLine.Matched[origin] = false
				cmdLine.Matched[len(cmdLine.Arguments)] = false
			}
		}
	}
}

// Available returns true if at least one argument is available.
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

// HasKey returns true if key is in Keys.
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

// HasValue returns true if value is in Values.
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

// HasValue returns true if origin is in Origins.
func (args *Arguments) HasOrigin(origin int) bool {
	if args != nil {
		for _, o := range args.Origins {
			if o == origin {
				return true
			}
		}
	}
	return false
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
					args.Origins = append(args.Origins, i)
				}
			}
		}
	}
	return args
}

func (cmdLine *CommandLine) searchPairsWithSpace(searchedKeys []string) *Arguments {
	var args *Arguments
	length, allMatched := len(cmdLine.Arguments), true
	for i := 0; i < length; i++ {
		if !cmdLine.Matched[i] {
			argument := cmdLine.Arguments[i]
			for _, searchedKey := range searchedKeys {
				if len(argument) > len(searchedKey) {
					value, ok := cmdLine.Delimiter.argumentValue(argument, searchedKey)
					if ok {
						if args == nil {
							args = new(Arguments)
						}
						args.Keys = append(args.Keys, searchedKey)
						args.Values = append(args.Values, value)
						args.Origins = append(args.Origins, i)
						cmdLine.Matched[i] = true
						break
					}
				} else if argument == searchedKey {
					value, iNxt := "", i+1
					if args == nil {
						args = new(Arguments)
					}
					args.Origins = append(args.Origins, i)
					cmdLine.Matched[i] = true
					if iNxt < length && !cmdLine.Matched[iNxt] {
						value = cmdLine.Arguments[iNxt]
						cmdLine.Matched[iNxt] = true
						i = iNxt
					}
					args.Keys = append(args.Keys, searchedKey)
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

func (cmdLine *CommandLine) searchPairsWithoutSpace(searchedKeys []string) *Arguments {
	var args *Arguments
	allMatched := true
	for i, argument := range cmdLine.Arguments {
		if !cmdLine.Matched[i] {
			for _, searchedKey := range searchedKeys {
				if len(argument) > len(searchedKey) {
					value, ok := cmdLine.Delimiter.argumentValue(argument, searchedKey)
					if ok {
						if args == nil {
							args = new(Arguments)
						}
						args.Keys = append(args.Keys, searchedKey)
						args.Values = append(args.Values, value)
						args.Origins = append(args.Origins, i)
						cmdLine.Matched[i] = true
						break
					}
				} else if cmdLine.Delimiter.HasEmptySeparator && argument == searchedKey {
					if args == nil {
						args = new(Arguments)
					}
					args.Keys = append(args.Keys, searchedKey)
					args.Values = append(args.Values, "")
					args.Origins = append(args.Origins, i)
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

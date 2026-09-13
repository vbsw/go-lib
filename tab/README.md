# tab

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/go-lib/tab.svg)](https://pkg.go.dev/github.com/vbsw/go-lib/tab)

## About
Package tab provides a parser for a simple, tab-indented data format. The format represents hierarchical structures using leading tab characters. Package tab is published on <https://github.com/vbsw/go-lib>.

## Copyright
Copyright 2026, Vitali Baumtrok (vbsw@mailbox.org).

tab is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

tab is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Usage

Code:

	package main

	import (
		"fmt"
		"github.com/vbsw/go-lib/tab"
	)

	func main() {
		var parser tab.ParserS
		data := "entities\n\tentity a001\n\t\tsize 500"

		parser.ParseLastLine = true
		for parser.Next(data) {
			switch parser.Key(data) {
			case "entities":
				if parser.Indent == 0 {
					fmt.Println("header OK        line", parser.LineNumber)
				}
			case "entity":
				if parser.Indent == 1 {
					fmt.Println("  entity OK      line", parser.LineNumber)
					fmt.Println("    name", parser.Value(data), "   line", parser.LineNumber)
				}
			case "size":
				if parser.Indent == 2 {
					fmt.Println("    size", parser.Value(data), "    line", parser.LineNumber)
				}
			}
		}
	}

Output:

	header OK        line 1
	  entity OK      line 2
	    name a001    line 2
	    size 500     line 3

## References
- https://go.dev/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git

# cl

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/go-lib/cl.svg)](https://pkg.go.dev/github.com/vbsw/go-lib/cl) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/go-lib/cl)](https://goreportcard.com/report/github.com/vbsw/go-lib/cl)

## About
Package cl parses command line arguments. It is published on <https://github.com/vbsw/go-lib>.

## Copyright
Copyright 2025, 2026, Vitali Baumtrok (vbsw@mailbox.org).

cl is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

cl is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Usage

### Example A

	package main

	import (
		"fmt"
		"github.com/vbsw/go-lib/cl"
		"os"
	)

	func main() {
		osCmdLine := cl.New(os.Args[1:], nil)

		if osCmdLine.Match("--help", "-h").Available() {
			fmt.Println("USAGE")
			fmt.Println("    --help         prints help")
			fmt.Println("    --version      prints version")

		} else if osCmdLine.Match("--version", "-v").Available() {
			fmt.Println("version 1.0.0")

		} else {
			unmatched := osCmdLine.Unmatched()

			if unmatched.Count() > 1 {
				fmt.Println("ERROR too many arguments")

			} else if unmatched.Count() == 1 {
				fmt.Printf("ERROR unknown argument \"%s\"\n", unmatched.Keys[0])
			}
		}
	}

### Example B

	package main

	import (
		"fmt"
		"github.com/vbsw/go-lib/cl"
		"os"
	)

	const defaultStart = "0"

	func main() {
		osCmdLine := cl.New(os.Args[1:], cl.NewDelimiter("="))

		start := osCmdLine.MatchDelimited("start").ValueAt(0, defaultStart)
		end := osCmdLine.MatchDelimited("end").ValueAt(0, start)

		fmt.Println("processing from", start, "to", end)
	}

## References
- https://go.dev/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git

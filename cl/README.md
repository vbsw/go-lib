# cl

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/go-lib/cl.svg)](https://pkg.go.dev/github.com/vbsw/go-lib/cl) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/go-lib/cl)](https://goreportcard.com/report/github.com/vbsw/go-lib/cl) [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## About
cl is a package for Go to parse command line arguments. It is published on <https://github.com/vbsw/go-lib/cl>.

## Copyright
Copyright 2025, Vitali Baumtrok (vbsw@mailbox.org).

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
		osCmdLine := cl.New(os.Args[1:])

		if osCmdLine.Search("--help", "-h").Available() {
			fmt.Println("USAGE")
			fmt.Println("    --help         prints help")
			fmt.Println("    --version      prints version")

		} else if osCmdLine.Search("--version", "-v").Available() {
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

	func main() {
		start := "0"
		end := "0"
		osCmdLine := cl.New(os.Args[1:])
		osCmdLine.Delimiter = cl.NewDelimiter("=")

		startArg := osCmdLine.SearchByDelimiter("start")
		endArg := osCmdLine.SearchByDelimiter("end")

		if startArg.Available() {
			start = startArg.Values[0]
			end = start
		}
		if endArg.Available() {
			end = endArg.Values[0]
		}
		fmt.Println("processing from", start, "to", end)
	}

## References
- https://go.dev/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git

# match

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/go-experimental/match.svg)](https://pkg.go.dev/github.com/vbsw/go-experimental/match) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/go-experimental/match)](https://goreportcard.com/report/github.com/vbsw/go-experimental/match)

## About
Package match provides simple wildcard string matching. Package match is published on <https://github.com/vbsw/go-experimental/match>.

## Copyright
Copyright 2026, Vitali Baumtrok (vbsw@mailbox.org).

cl is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

cl is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Compile
This package needs Cgo to compile and Cgo needs a C compiler.

**Linux**  
For Cgo install GCC, or configure another compiler like clang (see <https://stackoverflow.com/questions/44856124/can-i-change-default-compiler-used-by-cgo>).

**Windows**  
For Cgo install tdm-gcc (<https://jmeubank.github.io/tdm-gcc/>), or some other Go ABI compatible compiler like MinGW-w64.

## Execute

run tests

	go test

run benchmarks

	go test -bench=.

## References
- https://go.dev/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git

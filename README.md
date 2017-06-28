# Tomlvar

[![Build Status](https://travis-ci.org/dyson/envvar.svg?branch=master)](https://travis-ci.org/dyson/tomlvar)
[![Coverage Status](https://coveralls.io/repos/github/dyson/tomlvar/badge.svg?branch=master)](https://coveralls.io/github/dyson/tomlvar?branch=master)
[![Code Climate](https://codeclimate.com/github/dyson/tomlvar/badges/gpa.svg)](https://codeclimate.com/github/dyson/tomlvar)
[![Go Report Card](https://goreportcard.com/badge/github.com/dyson/tomlvar)](https://goreportcard.com/report/github.com/dyson/tomlvar)

[![GoDoc](https://godoc.org/github.com/dyson/tomlvar?status.svg)](http://godoc.org/github.com/dyson/tomlvar)
[![license](https://img.shields.io/github/license/dyson/tomlvar.svg)](https://github.com/dyson/tomlvar/blob/master/LICENSE)

Go toml config parsing in the style of flag.

Tomlvar is a fork and modification to the official Go flag package (https://golang.org/pkg/flag/). It has retained everything from flag that makes sense in the context of parsing toml config variables and removed everything else.

General use of the two packages are the same with the notable exception of:
 - Usage information for toml variables are not included.
 - A toml config must be loaded for the set before Parse() is called. Load(), LoadFile() and LoadReader() are supported from [go-toml](https://github.com/pelletier/go-toml/).
 - Uses [go-toml](https://github.com/pelletier/go-toml/) for parsing and retrieving toml configs.

## Documentation
https://godoc.org/github.com/dyson/tomlvar

## Installation
Using dep for dependency management (https://github.com/golang/dep):
```
dep ensure github.com/dyson/tomlvar
```

Using go get:
```
$ go get github.com/dyson/tomlvar
```
## Usage
Usage is essentially the same as the flag package. Here is an example program demonstrating tomlvar and the flag package being used together.

```
// example.go
package main

import (
	"flag"
	"fmt"

	"github.com/dyson/tomlvar"
)

type conf struct {
	a int
	b int
	c int
}

func main() {
	conf := &conf{
		a: 1,
		b: 1,
		c: 1,
	}

	// define flags and envvars
	flag.IntVar(&conf.a, "a", conf.a, "Value of a")
	tomlvar.IntVar(&conf.a, "letters.a", conf.a)

	flag.IntVar(&conf.b, "b", conf.b, "Value of b")
	tomlvar.IntVar(&conf.b, "letters.b", conf.b)
	
	flag.IntVar(&conf.c, "c", conf.c, "Value of c")
	tomlvar.IntVar(&conf.c, "letters.c", conf.c)

	// load toml string
	// could also use LoadFile(path string) to load file
	// or LoadReader(reader io.Reader) to load reader
	err := tomlvar.Load(`
[letters]
a = 2
b = 2
`)
	if err != nil {
		fmt.Println(err)
	}

	// parse in reverse precedence order
	// flags overwrite toml variables in this example
	tomlvar.Parse()
	flag.Parse()
	
	// print results
	fmt.Println("a set by flag precedence:", conf.a)
	fmt.Println("b set by toml var as no flag set:", conf.b) 
	fmt.Println("c set to default value as neither flag or toml var set it:", conf.c)
}
```

Running example:
```
$ go run example.go -a 3
a set by flag precedence: 3
b set by toml var as no flag set: 2
c set to default value as neither flag or toml var set it: 1
```

## Updates against flag
With tomlvar being so closely related to the flag package it makes sense to keep an eye on it's commits to see what bug fixes, improvements and features should be carried over to tomlvar.

Envvar was last checked against https://github.com/golang/go/tree/master/src/flag commit c65ceff125ded084c6f3b47f830050339e7cc74e.

If the above commit is not the latest commit to the flag package please submit an issue. This README should always reflect that tomlvar has been checked against the last flag commit.

## License
See LICENSE file.

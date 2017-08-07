# Tomlvar

[![Go project version](https://badge.fury.io/go/github.com%2Fdyson%2Ftomlvar.svg)](https://badge.fury.io/go/github.com%2Fdyson%2Ftomlvar)
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
### Basic
Usage is essentially the same as the flag package. Here is an example program demonstrating tomlvar and the flag package being used together.

```go
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
### Live reloading example
Here is an example with config reloading on SIGHUP
```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/dyson/tomlvar"
)

type config struct {
	mu sync.RWMutex
	a  int
}

func (c *config) Load() {
	a := 0
	tomlvar.TomlVars = tomlvar.NewTomlVarSet(os.Args[0], tomlvar.ExitOnError)

	tomlvar.IntVar(&a, "letters.a", a)

	path, err := filepath.Abs("./config.toml")
	if err == nil {
		if err = tomlvar.LoadFile(path); err == nil {
			tomlvar.Parse()
		}
	}
	if err != nil {
		fmt.Printf("skipping config file: %v\n", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.a = a
}

func (c *config) getA() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.a
}

func main() {
	fmt.Println("PID:", os.Getpid())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig)

	done := make(chan bool, 1)

	c := new(config)
	c.Load()

	go func() {
		for {
			sigReceived := <-sig
			switch sigReceived {
			case os.Interrupt:
				done <- true
			case syscall.SIGHUP:
				fmt.Println("reloading config")
				c.Load()
			}
		}
	}()

loop:
	for {
		select {
		case <-done:
			break loop
		default:
			fmt.Println("the value of a is:", c.getA())
			time.Sleep(time.Second)
		}
	}

	fmt.Println("exiting")
}
```
And the config.toml in the same directory:
```
[letters]
a = 1
```


Running example:
```
go run main.go
PID: 5427
the value of a is: 1
the value of a is: 1
...
```

Change the value of `a` in config.toml to 2 and kill -HUP 5427:
```
...
reloading config
the value of a is: 2
the value of a is: 2
^Cexiting
```

## Updates against flag
With tomlvar being so closely related to the flag package it makes sense to keep an eye on it's commits to see what bug fixes, improvements and features should be carried over to tomlvar.

Envvar was last checked against https://github.com/golang/go/tree/master/src/flag commit c65ceff125ded084c6f3b47f830050339e7cc74e.

If the above commit is not the latest commit to the flag package please submit an issue. This README should always reflect that tomlvar has been checked against the last flag commit.

## License
See LICENSE file.

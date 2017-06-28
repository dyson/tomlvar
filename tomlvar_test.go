// Copyright 2017 Dyson Simmons. All rights reserved.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tomlvar_test

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/dyson/tomlvar"
	"github.com/pelletier/go-toml"
)

func boolString(s string) string {
	if s == "0" {
		return "false"
	}
	return "true"
}

func TestEverything(t *testing.T) {
	ResetForTesting()

	Bool("test.bool", false)
	Int("test.int", 0)
	Int64("test.int64", 0)
	Uint("test.uint", 0)
	Uint64("test.uint64", 0)
	String("test.string", "0")
	Float64("test.float64", 0)
	Duration("test.duration", 0)

	m := make(map[string]*TomlVar)
	desired := "0"
	visitor := func(tv *TomlVar) {
		if len(tv.Path) > 5 && tv.Path[0:5] == "test." {
			m[tv.Path] = tv
			ok := false
			switch {
			case tv.Value.String() == desired:
				ok = true
			case tv.Path == "test.bool" && tv.Value.String() == boolString(desired):
				ok = true
			case tv.Path == "test.duration" && tv.Value.String() == desired+"s":
				ok = true
			}
			if !ok {
				t.Error("Visit: bad value", tv.Value.String(), "for", tv.Path)
			} else {
			}
		}
	}
	VisitAll(visitor)
	if len(m) != 8 {
		t.Error("VisitAll misses some defined envvars")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	m = make(map[string]*TomlVar)
	Visit(visitor)
	if len(m) != 0 {
		t.Errorf("Visit sees unset envvars")
		for k, v := range m {
			t.Log(k, *v)
		}
	}

	err := Load(`
[test]
bool = true
int = 1
int64 = 1
uint = 1
uint64 = 1
string = "1"
float64 = 1.0
duration = "1s"
`)
	if err != nil {
		t.Error(err)
	}

	// Now set all envvars
	Set("test.bool")
	Set("test.int")
	Set("test.int64")
	Set("test.uint")
	Set("test.uint64")
	Set("test.string")
	Set("test.float64")
	Set("test.duration")
	desired = "1"
	Visit(visitor)
	if len(m) != 8 {
		t.Error("Visit fails after set")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now test they're visited in sort order.
	var envVarNames []string
	Visit(func(tv *TomlVar) { envVarNames = append(envVarNames, tv.Path) })
	if !sort.StringsAreSorted(envVarNames) {
		t.Errorf("envvar names not sorted: %v", envVarNames)
	}
}

func TestGet(t *testing.T) {
	ResetForTesting()
	Bool("test.bool", true)
	Int("test.int", 1)
	Int64("test.int64", 2)
	Uint("test.uint", 3)
	Uint64("test.uint64", 4)
	String("test.string", "5")
	Float64("test.float64", 6)
	Duration("test.duration", 7)

	visitor := func(tv *TomlVar) {
		if len(tv.Path) > 5 && tv.Path[0:5] == "TEST_" {
			g, ok := tv.Value.(Getter)
			if !ok {
				t.Errorf("Visit: value does not satisfy Getter: %T", tv.Value)
				return
			}
			switch tv.Path {
			case "test.bool":
				ok = g.Get() == true
			case "test.int":
				ok = g.Get() == int(1)
			case "test.int64":
				ok = g.Get() == int64(2)
			case "test.uint":
				ok = g.Get() == uint(3)
			case "test.uint64":
				ok = g.Get() == uint64(4)
			case "test.string":
				ok = g.Get() == "5"
			case "test.float64":
				ok = g.Get() == float64(6)
			case "test.duration":
				ok = g.Get() == time.Duration(7)
			}
			if !ok {
				t.Errorf("Visit: bad value %T(%v) for %s", g.Get(), g.Get(), tv.Path)
			}
		}
	}
	VisitAll(visitor)
}

func testParse(tvs *TomlVarSet, t *testing.T) {
	if tvs.Parsed() {
		t.Error("tvs.Parse() = true before Parse")
	}
	boolTomlVar := tvs.Bool("test.bool", false)
	intTomlVar := tvs.Int("test.int", 0)
	int64TomlVar := tvs.Int64("test.int64", 0)
	uintTomlVar := tvs.Uint("test.uint", 0)
	uint64TomlVar := tvs.Uint64("test.uint64", 0)
	stringTomlVar := tvs.String("test.string", "0")
	float64TomlVar := tvs.Float64("test.float64", 0)
	durationTomlVar := tvs.Duration("test.duration", 5*time.Second)

	err := tvs.Load(`
[test]
bool = true
int = 22
int64 = 23
uint = 24
uint64 = 25
string = "hello"
float64 = 2718e28
duration = "2m"
`)
	if err != nil {
		t.Error(err)
	}

	if err := tvs.Parse(); err != nil {
		t.Fatal(err)
	}

	if !tvs.Parsed() {
		t.Error("tvs.Parse() = false after Parse")
	}
	if *boolTomlVar != true {
		t.Error("bool envvar should be true, is ", *boolTomlVar)
	}
	if *intTomlVar != 22 {
		t.Error("int envvar should be 22, is ", *intTomlVar)
	}
	if *int64TomlVar != 23 {
		t.Error("int64 envvar should be 23, is ", *int64TomlVar)
	}
	if *uintTomlVar != 24 {
		t.Error("uint envvar should be 24, is ", *uintTomlVar)
	}
	if *uint64TomlVar != 25 {
		t.Error("uint64 envvar should be 25, is ", *uint64TomlVar)
	}
	if *stringTomlVar != "hello" {
		t.Error("string envvar should be `hello`, is ", *stringTomlVar)
	}
	if *float64TomlVar != 2718e28 {
		t.Error("float64 envvar should be 2718e28, is ", *float64TomlVar)
	}
	if *durationTomlVar != 2*time.Minute {
		t.Error("duration envvar should be 2m, is ", *durationTomlVar)
	}
}

func TestTomlVarSetParse(t *testing.T) {
	testParse(NewTomlVarSet("test", ContinueOnError), t)
}

// Declare a user-defined envvar type.
type userVar []string

func (uv *userVar) String() string {
	return fmt.Sprint([]string(*uv))
}

func (uv *userVar) Set(path string, config *toml.Tree) error {
	for _, p := range strings.Split(path, ",") {
		v1 := config.Get(p)
		if v1 == nil {
			return nil
		}
		v2, ok := v1.(string)
		if !ok {
			return fmt.Errorf("can't convert \"%v\" (%T) to int", v1, v1)
		}
		*uv = userVar(append(*uv, v2))
	}
	return nil
}

func TestUserDefined(t *testing.T) {
	var tvs TomlVarSet
	tvs.Init("test", ContinueOnError)
	var uv userVar
	tvs.Var(&uv, "a,b,c")

	err := tvs.Load(`
a = "a"
b = "b"
c = "c"
`)
	if err != nil {
		t.Error(err)
	}

	if err := tvs.Parse(); err != nil {
		t.Fatal(err)
	}

	expect := "[a b c]"
	if uv.String() != expect {
		t.Errorf("expected value %q got %q", expect, uv.String())
	}
}

// Declare a user-defined boolean envvar type.
type boolTomlVar struct {
	count int
}

func (b *boolTomlVar) String() string {
	return fmt.Sprintf("%d", b.count)
}

func (b *boolTomlVar) Set(path string, config *toml.Tree) error {
	for _, p := range strings.Split(path, ",") {
		v1 := config.Get(p)
		if v1 == nil {
			return nil
		}
		v2, ok := v1.(bool)
		if !ok {
			return fmt.Errorf("can't convert \"%v\" (%T) to bool", v1, v1)
		}

		if bool(v2) == true {
			b.count++
		}
	}
	return nil
}

func TestUserDefinedBool(t *testing.T) {
	var tvs TomlVarSet
	tvs.Init("test", ContinueOnError)
	var b boolTomlVar
	tvs.Var(&b, "a,b,c")

	err := tvs.Load(`
a = true
b = false
c = true
notdefined = "something"
`)
	if err != nil {
		t.Error(err)
	}

	if err := tvs.Parse(); err != nil {
		t.Fatal(err)
	}

	if b.count != 2 {
		t.Errorf("want: %d; got: %d", 2, b.count)
	}
}

// Issue 19230 (from original flag package https://github.com/golang/go/): validate range of
// int and Uint TomlVar values.
func TestIntTomlVarOverflow(t *testing.T) {
	if strconv.IntSize != 32 {
		return
	}
	ResetForTesting()
	Int("i", 0)
	Uint("u", 0)

	err := Load(`
i = 2147483648
u = 4294967296
`)
	if err != nil {
		t.Error(err)
	}

	if err := Set("i"); err == nil {
		t.Error("unexpected success setting Int")
	}
	if err := Set("u"); err == nil {
		t.Error("unexpected success setting Uint")
	}
}

// Copyright 2017 Dyson Simmons. All rights reserved.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package tomlvar implements toml variable parsing.

Usage:

Parse a toml config using either Load, LoadFile or LoadReader.
	import "tomlvar"
	if err := tomlvar.LoadFile('config.toml'); err != nil {
		fmt.Println(err)
	}

Define toml variables using tomlvar.String(), Bool(), Int(), etc.

This declares an integer tomlvar, ENVVARNAME, stored in the pointer ip, with type *int.
	var ip = tomlvar.Int("ENVVARNAME", 1234)
If you like, you can bind the tomlvar to a variable using the Var() functions.
	var i int
	func init() {
		tomlvar.IntVar(&i, "ENVVARNAME", 1234)
	}
Or you can create custom tomlvars that satisfy the Value interface (with
pointer receivers) and couple them to toml variable parsing by
	tomlvar.Var(&tomlVarVal, "ENVVARNAME")
For such tomlvars, the default value is just the initial value of the variable.

After all tomlvars are defined, call
	tomlvar.Parse()
to parse the toml variables into the defined tomlvars.

Envvars may then be used directly. If you're using the tomlvars themselves,
they are all pointers; if you bind to variables, they're values.
	fmt.Println("ip has value ", *ip)
	fmt.Println("i has value ", i)

Integer tomlvars accept 1234, 0664, 0x1234 and may be negative.
Boolean tomlvars may be:
	1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False
Duration tomlvars accept any input valid for time.ParseDuration.

The default set of tomlvars is controlled by top-level functions.
The TomlVarSet type allows one to define	independent sets of tomlvars,
which facilitates their independent parsing. The methods of TomlVarSet
are	analogous to the top-level functions for the default	tomlvar set.
*/
package tomlvar

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/pelletier/go-toml"
)

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(bool)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to bool", v1, v1)
	}
	*b = boolValue(bool(v2))
	return nil
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(int64)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to int", v1, v1)
	}
	*i = intValue(int(v2))
	return nil
}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

// -- int64 Value
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(int64)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to int64", v1, v1)
	}
	*i = int64Value(v2)
	return nil
}

func (i *int64Value) Get() interface{} { return int64(*i) }

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// -- uint Value
type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(int64)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to uint", v1, v1)
	}
	*i = uintValue(uint(v2))
	return nil
}

func (i *uintValue) Get() interface{} { return uint(*i) }

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- uint64 Value
type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(int64)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to uint64", v1, v1)
	}
	*i = uint64Value(uint64(v2))
	return nil
}

func (i *uint64Value) Get() interface{} { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(string)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to string", v1, v1)
	}
	*s = stringValue(v2)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// -- float64 Value
type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(float64)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to float64", v1, v1)
	}
	*f = float64Value(v2)
	return nil
}

func (f *float64Value) Get() interface{} { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(path string, config *toml.Tree) error {
	v1 := config.Get(path)
	if v1 == nil {
		return nil
	}
	v2, ok := v1.(string)
	if !ok {
		return fmt.Errorf("can't convert \"%v\" (%T) to time.Duration", v1, v1)
	}
	v3, err := time.ParseDuration(v2)
	*d = durationValue(v3)
	return err
}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// Value is the interface to the dynamic value stored in a TomlVar.
// (The default value is represented as a string.)
//
// Set is called once for each TomlVar present.
// The tomlvar package may call the String method with a zero-valued receiver,
// such as a nil pointer.
type Value interface {
	String() string
	Set(string, *toml.Tree) error
}

// Getter is an interface that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
type Getter interface {
	Value
	Get() interface{}
}

// ErrorHandling defines how TomlVarSet.Parse behaves if the parse fails.
type ErrorHandling int

// These constants cause TomlVarSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError ErrorHandling = iota // return a descriptive error.
	ExitOnError                          // call os.Exit(2).
	PanicOnError                         // call panic with a descriptive error.
)

// A TomlVarSet represents a set of defined tomlVars. The zero value of a TomlVarSet
// has no name and has ContinueOnError error handling.
type TomlVarSet struct {
	name          string
	parsed        bool
	actual        map[string]*TomlVar
	formal        map[string]*TomlVar
	config        *toml.Tree
	errorHandling ErrorHandling
	output        io.Writer // nil means stderr; use out() accessor
}

// A TomlVar represents the state of a TomlVar.
type TomlVar struct {
	Path  string // path of toml variable
	Value Value  // value as set
}

// sortTomlVars returns the TomlVars as a slice in lexicographical sorted order.
func sortTomlVars(tomlVars map[string]*TomlVar) []*TomlVar {
	list := make(sort.StringSlice, len(tomlVars))
	i := 0
	for _, tv := range tomlVars {
		list[i] = tv.Path
		i++
	}
	list.Sort()
	result := make([]*TomlVar, len(list))
	for i, path := range list {
		result[i] = tomlVars[path]
	}
	return result
}

func (tvs *TomlVarSet) out() io.Writer {
	if tvs.output == nil {
		return os.Stderr
	}
	return tvs.output
}

// SetOutput sets the destination for error messages.
// If output is nil, os.Stderr is used.
func (tvs *TomlVarSet) SetOutput(output io.Writer) {
	tvs.output = output
}

// VisitAll visits the sets TomlVars in lexicographical order, calling
// fn for each. It visits all TomlVars, even those not set.
func (tvs *TomlVarSet) VisitAll(fn func(*TomlVar)) {
	for _, tomlVar := range sortTomlVars(tvs.formal) {
		fn(tomlVar)
	}
}

// VisitAll visits the default sets TomlVars in lexicographical order,
// calling fn for each. It visits TomlVars, even those not set.
func VisitAll(fn func(*TomlVar)) {
	TomlVars.VisitAll(fn)
}

// Visit visits the sets TomlVars in lexicographical order, calling fn for each.
// It visits only those TomlVars that have been set.
func (tvs *TomlVarSet) Visit(fn func(*TomlVar)) {
	for _, tomlVar := range sortTomlVars(tvs.actual) {
		fn(tomlVar)
	}
}

// Visit visits the default sets TomlVars in lexicographical order,
// calling fn for each. It visits only those TomlVars that have been set.
func Visit(fn func(*TomlVar)) {
	TomlVars.Visit(fn)
}

// Lookup returns the TomlVar structure of the named TomlVar,
// returning nil if none exists.
func (tvs *TomlVarSet) Lookup(path string) *TomlVar {
	return tvs.formal[path]
}

// Lookup returns the TomlVar structure of the named TomlVar,
// returning nil if none exists.
func Lookup(path string) *TomlVar {
	return TomlVars.formal[path]
}

// Set sets the value of the named TomlVar.
func (tvs *TomlVarSet) Set(path string) error {
	tomlVar, ok := tvs.formal[path]
	if !ok {
		return fmt.Errorf("no such tomlvar %v", path)
	}

	tomlVar.Value.Set(path, tvs.config)

	if tvs.actual == nil {
		tvs.actual = make(map[string]*TomlVar)
	}
	tvs.actual[path] = tomlVar
	return nil
}

// Set sets the value of the named TomlVar for the default set.
func Set(path string) error {
	return TomlVars.Set(path)
}

// NTomlVar returns the number of TomlVars that have been defined.
func (tvs *TomlVarSet) NTomlVar() int { return len(tvs.actual) }

// NTomlVar returns the number of TomlVars that have been defined.
func NTomlVar() int { return len(TomlVars.actual) }

// BoolVar defines a bool TomlVar with specified name, and default value.
// The argument p points to a bool variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) BoolVar(p *bool, path string, value bool) {
	tvs.Var(newBoolValue(value, p), path)
}

// BoolVar defines a bool TomlVar with specified name, and default value.
// The argument p points to a bool variable in which to store the value of the TomlVar.
func BoolVar(p *bool, path string, value bool) {
	TomlVars.Var(newBoolValue(value, p), path)
}

// Bool defines a bool TomlVar with specified name, and default value.
// The return value is the address of a bool variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) Bool(path string, value bool) *bool {
	p := new(bool)
	tvs.BoolVar(p, path, value)
	return p
}

// Bool defines a bool TomlVar with specified name, and default value.
// The return value is the address of a bool variable that stores the value of the TomlVar.
func Bool(path string, value bool) *bool {
	return TomlVars.Bool(path, value)
}

// IntVar defines an int TomlVar with specified name, and default value.
// The argument p points to an int variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) IntVar(p *int, path string, value int) {
	tvs.Var(newIntValue(value, p), path)
}

// IntVar defines an int TomlVar with specified name, and default value.
// The argument p points to an int variable in which to store the value of the TomlVar.
func IntVar(p *int, path string, value int) {
	TomlVars.Var(newIntValue(value, p), path)
}

// Int defines an int TomlVar with specified name, and default value.
// The return value is the address of an int variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) Int(path string, value int) *int {
	p := new(int)
	tvs.IntVar(p, path, value)
	return p
}

// Int defines an int TomlVar with specified name, and default value.
// The return value is the address of an int variable that stores the value of the TomlVar.
func Int(path string, value int) *int {
	return TomlVars.Int(path, value)
}

// Int64Var defines an int64 TomlVar with specified name, and default value.
// The argument p points to an int64 variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) Int64Var(p *int64, path string, value int64) {
	tvs.Var(newInt64Value(value, p), path)
}

// Int64Var defines an int64 TomlVar with specified name, and default value.
// The argument p points to an int64 variable in which to store the value of the TomlVar.
func Int64Var(p *int64, name string, value int64) {
	TomlVars.Var(newInt64Value(value, p), name)
}

// Int64 defines an int64 TomlVar with specified name, and default value.
// The return value is the address of an int64 variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) Int64(path string, value int64) *int64 {
	p := new(int64)
	tvs.Int64Var(p, path, value)
	return p
}

// Int64 defines an int64 TomlVar with specified name, and default value.
// The return value is the address of an int64 variable that stores the value of the TomlVar.
func Int64(path string, value int64) *int64 {
	return TomlVars.Int64(path, value)
}

// UintVar defines a uint TomlVar with specified name, and default value.
// The argument p points to a uint variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) UintVar(p *uint, path string, value uint) {
	tvs.Var(newUintValue(value, p), path)
}

// UintVar defines a uint TomlVar with specified name, and default value.
// The argument p points to a uint  variable in which to store the value of the TomlVar.
func UintVar(p *uint, path string, value uint) {
	TomlVars.Var(newUintValue(value, p), path)
}

// Uint defines a uint TomlVar with specified name, and default value.
// The return value is the address of a uint  variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) Uint(path string, value uint) *uint {
	p := new(uint)
	tvs.UintVar(p, path, value)
	return p
}

// Uint defines a uint TomlVar with specified name, and default value.
// The return value is the address of a uint  variable that stores the value of the TomlVar.
func Uint(path string, value uint) *uint {
	return TomlVars.Uint(path, value)
}

// Uint64Var defines a uint64 TomlVar with specified name, and default value.
// The argument p points to a uint64 variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) Uint64Var(p *uint64, path string, value uint64) {
	tvs.Var(newUint64Value(value, p), path)
}

// Uint64Var defines a uint64 TomlVar with specified name, and default value.
// The argument p points to a uint64 variable in which to store the value of the TomlVar.
func Uint64Var(p *uint64, path string, value uint64) {
	TomlVars.Var(newUint64Value(value, p), path)
}

// Uint64 defines a uint64 TomlVar with specified name, and default value.
// The return value is the address of a uint64 variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) Uint64(path string, value uint64) *uint64 {
	p := new(uint64)
	tvs.Uint64Var(p, path, value)
	return p
}

// Uint64 defines a uint64 TomlVar with specified name, and default value.
// The return value is the address of a uint64 variable that stores the value of the TomlVar.
func Uint64(path string, value uint64) *uint64 {
	return TomlVars.Uint64(path, value)
}

// StringVar defines a string TomlVar with specified name, and default value.
// The argument p points to a string variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) StringVar(p *string, path string, value string) {
	tvs.Var(newStringValue(value, p), path)
}

// StringVar defines a string TomlVar with specified name, and default value.
// The argument p points to a string variable in which to store the value of the TomlVar.
func StringVar(p *string, path string, value string) {
	TomlVars.Var(newStringValue(value, p), path)
}

// String defines a string TomlVar with specified name, and default value.
// The return value is the address of a string variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) String(path string, value string) *string {
	p := new(string)
	tvs.StringVar(p, path, value)
	return p
}

// String defines a string TomlVar with specified name, and default value.
// The return value is the address of a string variable that stores the value of the TomlVar.
func String(path string, value string) *string {
	return TomlVars.String(path, value)
}

// Float64Var defines a float64 TomlVar with specified name, and default value.
// The argument p points to a float64 variable in which to store the value of the TomlVar.
func (tvs *TomlVarSet) Float64Var(p *float64, path string, value float64) {
	tvs.Var(newFloat64Value(value, p), path)
}

// Float64Var defines a float64 TomlVar with specified name, and default value.
// The argument p points to a float64 variable in which to store the value of the TomlVar.
func Float64Var(p *float64, path string, value float64) {
	TomlVars.Var(newFloat64Value(value, p), path)
}

// Float64 defines a float64 TomlVar with specified name, and default value.
// The return value is the address of a float64 variable that stores the value of the TomlVar.
func (tvs *TomlVarSet) Float64(path string, value float64) *float64 {
	p := new(float64)
	tvs.Float64Var(p, path, value)
	return p
}

// Float64 defines a float64 TomlVar with specified name, and default value.
// The return value is the address of a float64 variable that stores the value of the TomlVar.
func Float64(path string, value float64) *float64 {
	return TomlVars.Float64(path, value)
}

// DurationVar defines a time.Duration TomlVar with specified name, and default value.
// The argument p points to a time.Duration variable in which to store the value of the TomlVar.
// The TomlVar accepts a value acceptable to time.ParseDuration.
func (tvs *TomlVarSet) DurationVar(p *time.Duration, path string, value time.Duration) {
	tvs.Var(newDurationValue(value, p), path)
}

// DurationVar defines a time.Duration TomlVar with specified name, and default value.
// The argument p points to a time.Duration variable in which to store the value of the TomlVar.
// The TomlVar accepts a value acceptable to time.ParseDuration.
func DurationVar(p *time.Duration, path string, value time.Duration) {
	TomlVars.Var(newDurationValue(value, p), path)
}

// Duration defines a time.Duration TomlVar with specified name, and default value.
// The return value is the address of a time.Duration variable that stores the value of the TomlVar.
// The TomlVar accepts a value acceptable to time.ParseDuration.
func (tvs *TomlVarSet) Duration(path string, value time.Duration) *time.Duration {
	p := new(time.Duration)
	tvs.DurationVar(p, path, value)
	return p
}

// Duration defines a time.Duration TomlVar with specified name, and default value.
// The return value is the address of a time.Duration variable that stores the value of the TomlVar.
// The TomlVar accepts a value acceptable to time.ParseDuration.
func Duration(path string, value time.Duration) *time.Duration {
	return TomlVars.Duration(path, value)
}

// Var defines a TomlVar with the specified name. The type and value of the TomlVar
// are represented by the first argument, of type Value, which typically holds a
// user-defined implementation of Value. For instance, the caller could create a
// TomlVar that turns a comma-separated string into a slice of strings by giving
// the slice the methods of Value; in particular, Set would decompose the
// comma-separated string into the slice.
func (tvs *TomlVarSet) Var(value Value, path string) {
	tomlVar := &TomlVar{path, value}
	_, alreadythere := tvs.formal[path]
	if alreadythere {
		var msg string
		if tvs.name == "" {
			msg = fmt.Sprintf("TomlVar redefined: %s", path)
		} else {
			msg = fmt.Sprintf("%s sets TomlVar redefined: %s", tvs.name, path)
		}
		fmt.Fprintln(tvs.out(), msg)
		panic(msg) // happens only if toml vars are declared with identical names
	}
	if tvs.formal == nil {
		tvs.formal = make(map[string]*TomlVar)
	}
	tvs.formal[path] = tomlVar
}

// Var defines a toml var with the specified name. The type and
// value of the toml var are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a toml var that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func Var(value Value, name string) {
	TomlVars.Var(value, name)
}

// failf prints to standard error a formatted error and returns the error.
func (tvs *TomlVarSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(tvs.out(), err)
	return err
}

// parseOne parses one toml var. It reports whether a toml var was seen.
func (tvs *TomlVarSet) parseOne(tomlVar *TomlVar) error {
	if err := tomlVar.Value.Set(tomlVar.Path, tvs.config); err != nil {
		return tvs.failf("invalid value for toml var %s: %v", tomlVar.Path, err)
	}
	if tvs.actual == nil {
		tvs.actual = make(map[string]*TomlVar)
	}
	tvs.actual[tomlVar.Path] = tomlVar
	return nil
}

// Parse parses all toml var definitions. Must be called after all toml vars in
// the TomlVarSet are defined and before toml vars are accessed by the program.
func (tvs *TomlVarSet) Parse() error {
	tvs.parsed = true

	for _, tomlVar := range tvs.formal {
		err := tvs.parseOne(tomlVar)
		if err != nil {
			switch tvs.errorHandling {
			case ContinueOnError:
				return err
			case ExitOnError:
				os.Exit(2)
			case PanicOnError:
				panic(err)
			}
		}
	}
	return nil
}

// Parsed reports whether tvs.Parse has been called.
func (tvs *TomlVarSet) Parsed() bool {
	return tvs.parsed
}

// Parse parses the toml vars from os.Environ().  Must be called
// after all toml vars are defined and before toml vars are accessed by the program.
func Parse() {
	TomlVars.Parse()
}

// Parsed reports whether the toml vars have been parsed.
func Parsed() bool {
	return TomlVars.Parsed()
}

// TomlVars is the default set of toml vars, parsed from os.Environ().
// The top-level functions such as BoolVar, Arg, and so on are wrappers for the
// methods of TomlVars.
var TomlVars = NewTomlVarSet(os.Args[0], ExitOnError)

// NewTomlVarSet returns a new, empty toml var set with the specified name and
// error handling property.
func NewTomlVarSet(name string, errorHandling ErrorHandling) *TomlVarSet {
	tvs := &TomlVarSet{
		name:          name,
		errorHandling: errorHandling,
	}
	return tvs
}

// Init sets the name and error handling property for a toml var set.
// By default, the zero TomlVarSet uses an empty name and the
// ContinueOnError error handling policy.
func (tvs *TomlVarSet) Init(name string, errorHandling ErrorHandling) {
	tvs.name = name
	tvs.errorHandling = errorHandling
}

// LoadReader creates a config Tree from any io.Reader.
func (tvs *TomlVarSet) LoadReader(reader io.Reader) error {
	var err error
	tvs.config, err = toml.LoadReader(reader)
	return err
}

// LoadReader creates a config Tree from any io.Reader.
func LoadReader(reader io.Reader) error {
	return TomlVars.LoadReader(reader)
}

// Load creates a config Tree from a toml string.
func (tvs *TomlVarSet) Load(content string) error {
	var err error
	tvs.config, err = toml.Load(content)
	return err
}

// Load creates a config Tree from a toml string.
func Load(content string) error {
	return TomlVars.Load(content)
}

// LoadFile creates a config Tree from a toml file.
func (tvs *TomlVarSet) LoadFile(path string) error {
	var err error
	tvs.config, err = toml.LoadFile(path)
	return err
}

// LoadFile creates a config Tree from a toml file.
func LoadFile(path string) error {
	return TomlVars.LoadFile(path)
}

// Config retrieves toml Tree.
func (tvs *TomlVarSet) Config() *toml.Tree {
	return tvs.config
}

// Config retrieves toml Tree.
func Config() *toml.Tree {
	return TomlVars.Config()
}

// Copyright 2017 Dyson Simmons. All rights reserved.

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tomlvar

import "os"

// Additional routines compiled into the package only during testing.

// ResetForTesting clears all default tomlvar state. After calling ResetForTesting,
// parse errors in tomlvar handling will not exit the program.
func ResetForTesting() {
	TomlVars = NewTomlVarSet(os.Args[0], ContinueOnError)
}

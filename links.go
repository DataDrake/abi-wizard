//
// Copyright 2019 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
	"os"
	"sort"
)

// libsFmt is the filename format for lib listings
const libsFmt = "abi%s_libs%s"

// synsFmt is the filename format for symbol listings
const symsFmt = "abi%s_symbols%s"

// Links models the linkage between libraries and symbols
type Links struct {
	libs []string
	syms Symbols
}

// Save writes a Links struct out to files as needed
func (l Links) Save(suffix, prefix string) {
	// ignore if empty list
	if len(l.libs) > 0 {
		// create the output file
		libs, err := os.Create(fmt.Sprintf(libsFmt, prefix, suffix))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create lib listing, reason: '%s'\n", err.Error())
			os.Exit(1)
		}
		// write out each library
		for i, lib := range l.libs {
			// skip duplicate libs
			if i > 0 && l.libs[i] == l.libs[i-1] {
				continue
			}
			fmt.Fprintln(libs, lib)
		}
		libs.Close()
	}
	// ignore if empty list
	if len(l.syms) > 0 {
		// create the output file
		syms, err := os.Create(fmt.Sprintf(symsFmt, prefix, suffix))
		if err != nil {
			panic(err.Error())
		}
		// write the symbols out
		l.syms.Print(syms)
		syms.Close()
	}
}

// Sort reorders the lists inside the Links struct
func (l Links) Sort() {
	sort.Strings(l.libs)
	sort.Sort(l.syms)
}

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
	"debug/elf"
	"fmt"
	"io"
)

// Symbols are lists of Symbols
type Symbols []elf.ImportedSymbol

// Len gives the length of Symbols for sorting
func (ss Symbols) Len() int {
	return len(ss)
}

// Less compares Symbols for sorting
func (ss Symbols) Less(i, j int) bool {
	// by library name
	if ss[i].Library < ss[j].Library {
		return true
	}
	if ss[i].Library > ss[j].Library {
		return false
	}
	/*
	       // by Version
	   	if ss[i].Version < ss[j].Version {
	   		return true
	   	}
	   	if ss[i].Version > ss[j].Version {
	   		return false
	   	}
	*/
	// by symbol name
	if ss[i].Name < ss[j].Name {
		return true
	}
	return false
}

// Swap switches entries for sorting
func (ss Symbols) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// Print writes smbols to a file
func (ss Symbols) Print(w io.Writer) {
	for i, symbol := range ss {
		// skip duplicates
		if i > 0 && ss[i-1].Name == ss[i].Name && ss[i-1].Library == ss[i].Library {
			continue
		}
		// replace missing library with UNKNOWN
		if len(symbol.Library) == 0 {
			symbol.Library = "UNKNOWN"
		}
		// write entry to file
		fmt.Fprintf(w, "%s:%s\n", symbol.Library, symbol.Name)
	}
}

// convertDynamic turns Exported symbols into Imported Symbols for simplicity
func convertDynamic(file string, ss []elf.Symbol) Symbols {
	syms := make(Symbols, 0)
	for _, s := range ss {
		sym := elf.ImportedSymbol{
			Name:    s.Name,
			Library: file,
		}
		syms = append(syms, sym)
	}
	return syms
}

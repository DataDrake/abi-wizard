//
// Copyright 2019-2021 Bryan T. Meyers <root@datadrake.com>
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

const (
	// libsFmt is the filename format for lib listings
	libsFmt = "abi%s_libs%s"
	// synsFmt is the filename format for symbol listings
	symsFmt = "abi%s_symbols%s"
)

// Links models the linkage between libraries and symbols
type Links struct {
	Libs map[string]int
	Syms map[string]Symbols
}

// NewLinks creates a new Links and its maps
func NewLinks() Links {
	return Links{
		Libs: make(map[string]int),
		Syms: make(map[string]Symbols),
	}
}

// Prune removes a set of related Links from another set of links
func (l Links) Prune(excludes Links) {
	for lib := range excludes.Libs {
		delete(l.Libs, lib)
		delete(l.Syms, lib)
	}
}

// Resolve will fix up any links where the library is unknown
func (l Links) Resolve(provided Links) {
	missingSymbols := l.Syms["UNKNOWN"]
	var unknown Symbols
	for _, missing := range missingSymbols {
		found := false
		for lib, symbols := range provided.Syms {
			for _, symbol := range symbols {
				if symbol == missing {
					l.Libs[lib]++
					l.Syms[lib] = append(l.Syms[lib], missing)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			unknown = append(unknown, missing)
		}
	}
	if len(unknown) == 0 {
		delete(l.Libs, "UNKNOWN")
		delete(l.Syms, "UNKNOWN")
		return
	}
	l.Libs["UNKNOWN"] = len(unknown)
	l.Syms["UNKNOWN"] = unknown
}

// Save writes a Links struct out to files as needed
func (l Links) Save(infix, suffix string) error {
	// ignore if empty list
	if len(l.Libs) == 0 {
		return nil
	}
	libs, err := os.Create(fmt.Sprintf(libsFmt, infix, suffix))
	if err != nil {
		return fmt.Errorf("failed to create lib listing, reason: '%s'", err)
	}
	defer libs.Close()
	syms, err := os.Create(fmt.Sprintf(symsFmt, infix, suffix))
	if err != nil {
		return fmt.Errorf("failed to create symbols listing, reason: '%s'", err)
	}
	defer syms.Close()
	var keys []string
	for key := range l.Libs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, lib := range keys {
		symbols := l.Syms[lib]
		fmt.Fprintln(libs, lib)
		sort.Sort(symbols)
		for i, symbol := range symbols {
			if i > 0 && symbols[i-1] == symbol {
				continue
			}
			fmt.Fprintf(syms, "%s:%s\n", lib, symbol)
		}
	}
	return nil
}

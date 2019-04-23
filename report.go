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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// machineNames maps machine type to a file suffix
var machineNames = map[elf.Machine]string{
	elf.EM_386:    "32",
	elf.EM_X86_64: "",
}

// Report contains one or more architecture descriptions
type Report struct {
	arches map[string]Arch
}

// Save writes a report to disk
func (r Report) Save() {
	for _, arch := range r.arches {
		arch.Save()
	}
}

// AddDir reads the contents of a directory and adds information to the report
func (r Report) AddDir(path string) {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, info := range infos {
		r.Add(filepath.Join(path, info.Name()))
	}
}

// Add the specified path to the report
func (r Report) Add(path string) {
	// try to stat the path
	info, err := os.Lstat(path)
	if err != nil {
		return
	}
	// use AddDir for direcctories
	if info.IsDir() {
		r.AddDir(path)
		return
	}
	// check for executable bit to ignore other file types
	mode := info.Mode()
	if mode&0111 == 0 {
		return
	}
	// ignore statically linked archives or debug symbols
	if strings.HasSuffix(info.Name(), ".la") || strings.HasSuffix(info.Name(), ".a") || strings.HasSuffix(info.Name(), ".debug") {
		return
	}
	// Read the ELF
	f, err := elf.Open(path)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "failed to open '%s', reason: '%s'\n", path, err.Error())
		//os.Exit(1)
		return
	}
	defer f.Close()
	// get the architecture entry or create if missing
	name := machineNames[f.Machine]
	arch := r.arches[name]
	if arch.suffix != name {
		arch = Arch{
			suffix: name,
		}
	}
	// process based on ELF type
	switch f.Type {
	case elf.ET_DYN:
		// Shared Object / DLL Libraries
		arch.provides.libs = append(arch.provides.libs, info.Name())
		symbols, err := f.DynamicSymbols()
		if err == nil {
			syms := convertDynamic(info.Name(), symbols)
			arch.provides.syms = append(arch.provides.syms, syms...)
		}
		// still need to process imports
		fallthrough
	case elf.ET_EXEC, elf.ET_REL:
		// Executables and relocatable binaries
		libs, err := f.ImportedLibraries()
		if err != nil {
			break
		}
		arch.uses.libs = append(arch.uses.libs, libs...)
		symbols, err := f.ImportedSymbols()
		if err == nil {
			arch.uses.syms = append(arch.uses.syms, symbols...)
		}
	default:
		return
	}
	// save changes
	r.arches[name] = arch
}

// Sort the architectures inside the Report
func (r Report) Sort() {
	for _, arch := range r.arches {
		arch.Sort()
	}
}

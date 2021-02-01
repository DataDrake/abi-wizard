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

package abi

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// machineNames maps machine type to a file suffix
var machineNames = map[elf.Machine]string{
	elf.EM_386:    "32",
	elf.EM_X86_64: "",
}

// machineTypes maps machine type to a file suffix
var machineTypes = map[string]elf.Machine{
	"32": elf.EM_386,
	"":   elf.EM_X86_64,
}

// machineLibs maps machine type to directories for its libs
var machineLibs = map[elf.Machine][]string{
	elf.EM_386: []string{
		"/lib32",
		"/usr/lib32",
	},
	elf.EM_X86_64: []string{
		"/lib64",
		"/usr/lib64",
	},
}

// Report contains one or more architecture descriptions
type Report map[string]Arch

// Resolve missing libraries
func (r Report) Resolve() (missing []string, err error) {
	for _, arch := range r {
		missing = append(missing, arch.Resolve()...)
	}
	var unique []string
	sort.Strings(missing)
	for i, m := range missing {
		if i != 0 && missing[i-1] == m {
			continue
		}
		unique = append(unique, m)
	}
	missing = make([]string, 0)
	if len(unique) > 0 {
		r2 := make(Report)
		for arch := range r {
			archType := machineTypes[arch]
			for _, dir := range machineLibs[archType] {
				err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					name := info.Name()
					for _, u := range unique {
						if name == u {
							return r2.Add(path)
						}
					}
					return nil
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s\n", err)
				}
			}
		}
		for name, arch := range r {
			missing = append(missing, arch.ResolveMissing(r2[name])...)
		}
	}
	sort.Strings(missing)
	return
}

// Save writes a report to disk
func (r Report) Save() error {
	for _, arch := range r {
		if err := arch.Save(); err != nil {
			return err
		}
	}
	return nil
}

// Add the specified path to the report
func (r Report) Add(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// check for executable bit to ignore other file types
		if info.Mode()&0111 == 0 {
			return nil
		}
		// ignore statically linked archives or debug symbols
		if strings.HasSuffix(info.Name(), ".la") || strings.HasSuffix(info.Name(), ".a") || strings.HasSuffix(info.Name(), ".debug") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file '%s', reason: %s", path, err)
		}
		defer f.Close()
		return r.AddFile(f, info.Name())
	})
}

// AddFile adds the specified file to the report
func (r Report) AddFile(in io.ReaderAt, name string) error {
	// Read the ELF
	f, err := elf.NewFile(in)
	if err != nil {
		if strings.Contains(err.Error(), "bad magic number") {
			return nil
		}
		return fmt.Errorf("failed to open '%s', reason: '%s'", name, err)
	}
	defer f.Close()
	// get the architecture entry or create if missing
	archName := machineNames[f.Machine]
	arch, ok := r[archName]
	if !ok {
		arch = NewArch(archName)
	}
	// process based on ELF type
	switch f.Type {
	case elf.ET_DYN:
		// Shared Object / DLL Libraries
		symbols, err := f.DynamicSymbols()
		if err != nil {
			return err
		}
		dynName, err := f.DynString(elf.DT_SONAME)
		if err != nil {
			return err
		}
		if len(dynName) > 0 {
			name = dynName[0]
		}
		for _, symbol := range symbols {
			stBind := elf.ST_BIND(symbol.Info)
			if (stBind & elf.STB_WEAK) == elf.STB_WEAK {
				continue
			}
			if symbol.Section == elf.SHN_UNDEF {
				continue
			}
			arch.Provides.Libs[name]++
			arch.Provides.Syms[name] = append(arch.Provides.Syms[name], symbol.Name)
		}
		fallthrough
	case elf.ET_EXEC, elf.ET_REL:
		// Executables and relocatable binaries
		libs, err := f.ImportedLibraries()
		for _, lib := range libs {
			arch.Uses.Libs[lib]++
		}
		symbols, err := f.ImportedSymbols()
		if err != nil {
			return err
		}
		for _, symbol := range symbols {
			name := symbol.Library
			if len(name) == 0 {
				name = "UNKNOWN"
			}
			arch.Uses.Libs[name]++
			arch.Uses.Syms[name] = append(arch.Uses.Syms[name], symbol.Name)
		}
	default:
		return nil
	}
	// save changes
	r[archName] = arch
	return nil
}

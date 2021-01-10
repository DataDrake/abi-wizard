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
	"debug/elf"
	"fmt"
	"io"
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
type Report map[string]Arch

// Resolve missing libraries
func (r Report) Resolve() {
	for _, arch := range r {
		arch.Resolve()
	}
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

// AddDir reads the contents of a directory and adds information to the report
func (r Report) AddDir(path string) error {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory '%s', reason: %s", path, err)
	}
	for _, info := range infos {
		if err = r.Add(filepath.Join(path, info.Name())); err != nil {
			return err
		}
	}
	return nil
}

// Add the specified path to the report
func (r Report) Add(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("could not stat file '%s', reason: %s", path, err)
	}
	if info.IsDir() {
		return r.AddDir(path)
	}
	if (info.Mode() & os.ModeSymlink) == os.ModeSymlink {
		return nil
	}
	// check for executable bit to ignore other file types
	if mode := info.Mode(); mode&0111 == 0 {
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

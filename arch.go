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

// Arch bundles all of the symbols and libraries for a given architecture
type Arch struct {
	Suffix   string
	Provides Links
	Uses     Links
}

// NewArch returns a new empty Architecture
func NewArch(suffix string) Arch {
	return Arch{
		Suffix:   suffix,
		Provides: NewLinks(),
		Uses:     NewLinks(),
	}
}

// Resolve unknown libraries
func (a *Arch) Resolve() []string {
	return a.Uses.Resolve(a.Provides)
}

// ResolveMissing unknown libraries from a search of the missing libraries
func (a *Arch) ResolveMissing(a2 Arch) []string {
	return a.Uses.Resolve(a2.Provides)
}

// Save writes an architecture to disk
func (a Arch) Save() error {
	a.Uses.Prune(a.Provides)
	if err := a.Provides.Save("", a.Suffix); err != nil {
		return err
	}
	return a.Uses.Save("_used", a.Suffix)
}

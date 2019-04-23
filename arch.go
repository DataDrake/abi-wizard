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

// Arch bundles all of the symbols and libraries for a given architecture
type Arch struct {
	suffix   string
	provides Links
	uses     Links
}

// Save writes an architecture to disk
func (a Arch) Save() {
	a.provides.Save(a.suffix, "")
	a.uses.Save(a.suffix, "_used")
}

// Sort reorders the entries inside the struct
func (a Arch) Sort() {
	a.provides.Sort()
	a.uses.Sort()
}

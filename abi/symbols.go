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

// Symbols are lists of Symbols
type Symbols []string

// Len gives the length of Symbols for sorting
func (ss Symbols) Len() int {
	return len(ss)
}

// Less compares Symbols for sorting
func (ss Symbols) Less(i, j int) bool {
	return ss[i] < ss[j]
}

// Swap switches entries for sorting
func (ss Symbols) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

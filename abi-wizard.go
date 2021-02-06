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
	"github.com/DataDrake/abi-wizard/abi"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: abi-wizard <file/folder>")
		os.Exit(1)
	}
	r := make(abi.Report)
	if err := r.Add(os.Args[1], os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	missing, err := r.Resolve()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	for _, lib := range missing {
		fmt.Fprintf(os.Stderr, "Missing library: %s\n", lib)
	}
	if err = r.Save("."); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

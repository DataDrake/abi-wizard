# abi-wizard
Tool for generating ABI reports for libraries and binaries

[![Go Report Card](https://goreportcard.com/badge/github.com/DataDrake/abi-wizard)](https://goreportcard.com/report/github.com/DataDrake/abi-wizard) [![license](https://img.shields.io/github/license/DataDrake/abi-wizard.svg)]() 

## Motivation

As a package maintainer, it's a challenging task to keep track of binary dependencies and ABI changes. This tool generates reports that can be used to inform maintainers of ABI and dependency changes.

## Goals
 * Fast scan times
 * Be completely distro agnostic
 * A+ Rating on [Report Card](https://goreportcard.com/report/github.com/DataDrake/abi-wizard)
 
## Installation

1. Clone repo and enter its
2. `make`
3. `sudo make install`

## Usage

abi-wizard <file/path>

where file/path is a path to a location to scan.

## Output

Up to 4 files per device architecture may be generated:

| File             | Purpose                                            |
| ---------------- | -------------------------------------------------- |
| abi_libs         | List of ELF files provided by the search path      |
| abi_symbols      | List of symbols exported by the detected ELF files |
| abi_used_libs    | List of libs imported by the detected ELF files    |
| abi_used_symbols | List of symbols imported by the detected ELF files |

A suffix will be added to the output files to signify the architecture:

| Architecture | Suffix |
| ------------ | ------ |
| x86          | 32     |
| x86_64       | N/A    |

Currently only 32-bit and 64-bit x86 architectures are supported

## License
 
Copyright 2019-2021 Bryan T. Meyers <root@datadrake.com>
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

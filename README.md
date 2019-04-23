# abi-wizard
Tool for generating ABI reports for libraries and binaries

[![Go Report Card](https://goreportcard.com/badge/github.com/DataDrake/abi-wizard)](https://goreportcard.com/report/github.com/DataDrake/abi-wizard) [![license](https://img.shields.io/github/license/DataDrake/abi-wizard.svg)]() 

## Motivation

As a package maintainer, it's a challenging task to keep track of binary dependencies and ABI changes. This tool generatres reports that can be used to inform maintainers of ABI and dependency changes.

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

## License
 
Copyright 2019 Bryan T. Meyers <bmeyers@datadrake.com>
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 

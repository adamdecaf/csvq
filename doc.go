// Licensed to Adam Shannon under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"strings"
)

func help() {
	fmt.Printf(strings.TrimSpace(fmt.Sprintf(`
csvq is a tool for parsing and transforming comma separated values (CSV) files. This tool allows
for flexible transforms of data files.

COMMAND LINE

   -d <rune>             Delimiter used to separate records. (Default: ',')
   -keep <string-list>   Column headers to keep in output. Order of kept headers is maintained in output.

   -headers              Show headers in output.

   -format <string>      Layout to print resulting records in. (Default: csv, Options: csv, tabs, table)

   -v                    Enable verbose logging
   -version              Print the version of csvq. (Example: %s)

EXAMPLES

Extract first_name and last_name columns (in that order). Sort results.
   csvq -keep first_name,last_name ~/Downloads/report.csv | sort -u

Change delimiter used in report.csv.
   csvq -d';' user_id,dob,email ~/Downloads/report.csv

Output CSV columns in a table.
   csvq -keep first_name,last_name -format table

Combine multiple files.
   csvq -keep user_id,email ~/Downloads/report1.csv ~/Downloads/report2.csv
`, Version)))

	// TODO(adam): support additional flags
	// // flagIgnoreCols = flag.String("ignore", "", "Column headers to remove from output")
	// // flagIndices    = flag.String("i", "", "Indicies to keep in output")
	// // flagNotIndices = flag.String("I", "", "Indicies to remove from output")

	fmt.Println("")

	// // TODO(adam): flag to sort output (useful w/ '-format table')
	// // IDEA: -sort.asc/-sort.desc 1 (index) or last_name (col name)

	// // TODO(adam): -exec flag to call external program on PATH
	// //  could we format command with $1, $2 or $first_name, $last_name values?
}

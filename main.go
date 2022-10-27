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
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/adamdecaf/csvq/internal/cli"
	"github.com/adamdecaf/csvq/internal/format"
)

var (
	flagDelimiter   = flag.String("d", ",", "Delimiter used to separate records")
	flagShowHeaders = flag.Bool("headers", false, "Print headers as first output line")

	flagKeepCols = flag.String("keep", "", "Column headers to keep in output")
	// flagIgnoreCols = flag.String("ignore", "", "Column headers to remove from output") // TODO(adam):
	// flagIndices    = flag.String("i", "", "Indicies to keep in output")
	// flagNotIndices = flag.String("I", "", "Indicies to remove from output")

	flagFormat = flag.String("format", "", "Format to output resulting records in") // TODO(adam): csv, tabs, table

	flagVerbose = flag.Bool("v", false, "Enable verbose logging")
	flagVersion = flag.Bool("version", false, "Print the version of csvq")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("csvq %s", Version)
		return
	}

	files, err := cli.OpenPaths(flag.Args())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := cli.FileOpts{
		Delimiter:   toRune(*flagDelimiter),
		ShowHeaders: *flagShowHeaders,
		KeepCols:    splitStringList(*flagKeepCols),
	}

	for i := range files {
		if *flagVerbose {
			fmt.Printf("Processing %s\n", files[i].Name())
		}

		output, err := cli.HandleFile(opts, files[i])
		if err != nil {
			fmt.Printf("ERROR with %s handler: %v", files[i], err)
			os.Exit(1)
		}
		err = format.WriteFile(os.Stdout, *flagFormat, output)
		if err != nil {
			fmt.Printf("ERROR writing output: %v", err)
		}
	}
}

func toRune(delim string) rune {
	if utf8.RuneCountInString(delim) != 1 {
		return ',' // delim is invalid
	}
	return rune(delim[0])
}

func splitStringList(input string) []string {
	ss := strings.Split(input, ",")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return ss
}

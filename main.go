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
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/adamdecaf/csvq/internal/cli"
)

var (
	flagDelimiter   = flag.String("d", ",", "Delimiter used to separate records")
	flagShowHeaders = flag.Bool("headers", false, "Print headers as first output line")
	flagKeepCols    = flag.String("keep", "", "Column headers to keep in output")

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

	for i := range files {
		if *flagVerbose {
			fmt.Printf("Processing %s\n", files[i].Name())
		}

		rdr := csv.NewReader(files[i])
		rdr.Comma = toRune(*flagDelimiter)

		headerIndexes := make(map[int]string)

		lineNumber := 0
		for {
			lineNumber += 1

			cols, err := rdr.Read()
			if err != nil {
				if err == io.EOF {
					if *flagVerbose {
						fmt.Println("Finished processing")
					}
					break
				}
				fmt.Printf("ERROR reading line %d failed: %v", lineNumber, err)
				os.Exit(1)
			}

			if lineNumber == 1 && *flagShowHeaders {
				fmt.Println(strings.Join(cols, ", ")) // TODO(adam): text/tabwriter
			}

			if lineNumber == 1 {
				toKeep := strings.Split(*flagKeepCols, ",")

				for i := range cols {
					for j := range toKeep {
						if strings.EqualFold(strings.TrimSpace(cols[i]), strings.TrimSpace(toKeep[j])) {
							headerIndexes[i] = strings.TrimSpace(cols[i])
						}
					}
				}

				fmt.Printf("%#v\n", headerIndexes)
			} else {
				for i := range cols {
					if _, exists := headerIndexes[i]; exists {
						// TODO(adam): tabwriter
						fmt.Printf("  %s", cols[i])
					}
				}
				fmt.Printf("\n")
			}
		}
	}
}

func toRune(delim string) rune {
	if utf8.RuneCountInString(delim) != 1 {
		return ',' // delim is invalid
	}
	return rune(delim[0])
}

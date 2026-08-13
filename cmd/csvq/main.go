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
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/adamdecaf/csvq"
	"github.com/adamdecaf/csvq/internal/cli"
	"github.com/adamdecaf/csvq/internal/format"
)

var (
	flagDelimiter   = flag.String("d", ",", "Delimiter used to separate records")
	flagShowHeaders = flag.Bool("headers", false, "Print headers as first output line")

	flagKeepCols = flag.String("keep", "", "Column headers to keep in output. Order of kept headers is maintained in output.")

	_ = flag.String("sort.asc", "", "Comma-separated column headers to sort output by (ascending)")
	_ = flag.String("sort.dsc", "", "Comma-separated column headers to sort output by (descending)")

	flagFormat = flag.String("format", "", "Format to output resulting records in")

	flagVerbose = flag.Bool("v", false, "Enable verbose logging")
	flagVersion = flag.Bool("version", false, "Print the version of csvq")
)

//go:embed help.txt
var helpText string

func main() {
	flag.Usage = func() {
		fmt.Println(helpText)
	}
	flag.Parse()

	if *flagVersion {
		fmt.Printf("csvq %s", csvq.Version)
		return
	}

	files, err := cli.OpenPaths(flag.Args())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer files.Close()

	if len(files) == 0 {
		flag.Usage()
		return
	}

	opts := cli.FileOpts{
		Delimiter:   toRune(*flagDelimiter),
		ShowHeaders: *flagShowHeaders,
		KeepCols:    splitStringList(*flagKeepCols),
		SortKeys:    parseSortKeys(os.Args[1:]),
	}

	for i := range files {
		if *flagVerbose {
			fmt.Printf("Processing %s\n", files[i].Name())
		}

		output, err := cli.HandleFile(opts, files[i])
		if err != nil {
			fmt.Printf("ERROR with %s handler: %v", files[i].Name(), err)
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
	input = strings.TrimSpace(input)

	if input == "" {
		return nil
	}

	ss := strings.Split(input, ",")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return ss
}

// parseSortKeys walks argv so mixed -sort.asc / -sort.dsc flags keep invocation order.
func parseSortKeys(args []string) []cli.SortKey {
	var keys []cli.SortKey
	for i := 0; i < len(args); i++ {
		arg := args[i]
		val, desc, ok, skipNext := sortFlagValue(arg, args, i)
		if !ok {
			continue
		}
		if skipNext {
			i++
		}
		for _, name := range splitStringList(val) {
			keys = append(keys, cli.SortKey{Name: name, Desc: desc})
		}
	}
	return keys
}

func sortFlagValue(arg string, args []string, i int) (val string, desc bool, ok bool, skipNext bool) {
	name, inline, hasInline := strings.Cut(arg, "=")
	switch name {
	case "-sort.asc", "--sort.asc":
		desc = false
	case "-sort.dsc", "--sort.dsc":
		desc = true
	default:
		return "", false, false, false
	}
	if hasInline {
		return inline, desc, true, false
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
		return args[i+1], desc, true, true
	}
	return "", desc, true, false
}

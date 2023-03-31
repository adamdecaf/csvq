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

package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type FileOpts struct {
	Delimiter   rune
	ShowHeaders bool

	KeepCols []string
}

type File struct {
	Headers []string
	Lines   []Line
	Opts    FileOpts
}

type Line []string

func HandleFile(opts FileOpts, r io.Reader) (*File, error) {
	if r == nil {
		return nil, errors.New("missing file")
	}

	rdr := csv.NewReader(r)
	rdr.Comma = opts.Delimiter

	// Until we handle flagIndices or flagNotIndices we assume the first record contains headers
	headerIndexes := make(map[int]int) // index within input file -> index in output
	headers, err := rdr.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	output := File{
		Opts: opts,
	}

	toKeep := opts.KeepCols
	for headIdx := range headers {
		if len(toKeep) == 0 {
			headerIndexes[headIdx] = headIdx // keep all cols, so maintain order
			output.Headers = append(output.Headers, headers[headIdx])
		}

		for keepIdx := range toKeep {
			if strings.EqualFold(strings.TrimSpace(headers[headIdx]), strings.TrimSpace(toKeep[keepIdx])) {
				headerIndexes[headIdx] = keepIdx
				output.Headers = append(output.Headers, toKeep[keepIdx])
			}
		}
	}

	lineNumber := 1
	for {
		lineNumber += 1

		cols, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("line %d failed to parse: %v", lineNumber, err)
		}

		line := make(Line, len(toKeep))
		for i := range cols {
			if keepIdx, exists := headerIndexes[i]; exists {
				line[keepIdx] = cols[i] // format output according to -keep order
			}
		}

		output.Lines = append(output.Lines, line)
	}

	return &output, nil
}

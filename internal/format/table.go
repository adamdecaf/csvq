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

package format

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/adamdecaf/csvq/internal/cli"
)

func writeTable(ww io.Writer, file *cli.File) error {
	w := tabwriter.NewWriter(ww, 0, 0, 2, ' ', 0)

	if file.Opts.ShowHeaders {
		var line strings.Builder
		for idx, hdr := range file.Opts.KeepCols {
			if idx > 0 {
				line.WriteString("\t")
			}
			line.WriteString(hdr)
		}
		fmt.Fprintln(w, line.String())
	}

	var positions []int
	for _, name := range file.Opts.KeepCols {
		idx := slices.Index(file.Headers, name)
		positions = append(positions, idx)
	}

	for _, line := range file.Lines {
		var buf strings.Builder
		for col, posIdx := range positions {
			if col > 0 {
				buf.WriteString("\t")
			}
			buf.WriteString(line[posIdx])
		}
		fmt.Fprintln(w, buf.String())
	}

	return w.Flush()
}

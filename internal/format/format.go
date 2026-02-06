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
	"strings"

	"github.com/adamdecaf/csvq/internal/cli"
)

func WriteFile(w io.Writer, format string, file *cli.File) error {
	switch strings.ToLower(format) {
	case "", "csv":
		return writeCSV(w, file)
	case "tab":
		// TODO(adam):
	case "table":
		return writeTable(w, file)
	case "web":
		// TODO(adam):
	}
	return fmt.Errorf("unexpected %s format option", format)
}

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
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adamdecaf/csvq/internal/cli"

	"github.com/stretchr/testify/require"
)

func TestFormat_Table(t *testing.T) {
	opts := cli.FileOpts{
		Delimiter:   ',',
		ShowHeaders: true,
		KeepCols:    []string{"date", "score", "name"},
	}

	fd, err := os.Open(filepath.Join("..", "..", "testdata", "scores.csv"))
	require.NoError(t, err)

	t.Cleanup(func() { fd.Close() })

	file, err := cli.HandleFile(opts, fd)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = writeTable(&buf, file)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
date  score  name
a     12     2025
b     45     2026
c     32     2027
`)
	require.Equal(t, expected, strings.TrimSpace(buf.String()))
}

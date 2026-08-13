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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const scoresCSV = `name,score,date
a,12,2025
b,45,2026
c,32,2027
`

func TestHandleFile_Sort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		opts    FileOpts
		want    []Line
		wantErr string
	}{
		{
			name:  "asc by name without keep",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				SortKeys:  []SortKey{{Name: "name"}},
			},
			want: []Line{
				{"a", "12", "2025"},
				{"b", "45", "2026"},
				{"c", "32", "2027"},
			},
		},
		{
			name:  "desc by score without keep",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				SortKeys:  []SortKey{{Name: "score", Desc: true}},
			},
			want: []Line{
				{"b", "45", "2026"},
				{"c", "32", "2027"},
				{"a", "12", "2025"},
			},
		},
		{
			name:  "keep reorders columns then sort by score desc",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				KeepCols:  []string{"score", "name"},
				SortKeys:  []SortKey{{Name: "score", Desc: true}},
			},
			want: []Line{
				{"45", "b"},
				{"32", "c"},
				{"12", "a"},
			},
		},
		{
			name:  "keep in file order then sort by name asc",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				KeepCols:  []string{"name", "score"},
				SortKeys:  []SortKey{{Name: "name"}},
			},
			want: []Line{
				{"a", "12"},
				{"b", "45"},
				{"c", "32"},
			},
		},
		{
			name: "multi-column desc score then asc name",
			input: `name,score
c,10
a,20
b,10
`,
			opts: FileOpts{
				Delimiter: ',',
				SortKeys: []SortKey{
					{Name: "score", Desc: true},
					{Name: "name"},
				},
			},
			want: []Line{
				{"a", "20"},
				{"b", "10"},
				{"c", "10"},
			},
		},
		{
			name: "stable on ties",
			input: `name,group
first,x
second,x
third,y
`,
			opts: FileOpts{
				Delimiter: ',',
				SortKeys:  []SortKey{{Name: "group"}},
			},
			want: []Line{
				{"first", "x"},
				{"second", "x"},
				{"third", "y"},
			},
		},
		{
			name:  "unknown column",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				SortKeys:  []SortKey{{Name: "missing"}},
			},
			wantErr: `unknown sort column "missing"`,
		},
		{
			name:  "sort by dropped keep column",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				KeepCols:  []string{"name"},
				SortKeys:  []SortKey{{Name: "score", Desc: true}},
			},
			wantErr: `unknown sort column "score"`,
		},
		{
			name: "numeric not lexicographic",
			input: `name,score
a,9
b,45
c,12
`,
			opts: FileOpts{
				Delimiter: ',',
				SortKeys:  []SortKey{{Name: "score"}},
			},
			want: []Line{
				{"a", "9"},
				{"c", "12"},
				{"b", "45"},
			},
		},
		{
			name:  "case-insensitive column names",
			input: scoresCSV,
			opts: FileOpts{
				Delimiter: ',',
				KeepCols:  []string{"SCORE", "NAME"},
				SortKeys:  []SortKey{{Name: "Score", Desc: true}},
			},
			want: []Line{
				{"45", "b"},
				{"32", "c"},
				{"12", "a"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := HandleFile(tt.opts, strings.NewReader(tt.input))
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Lines)
		})
	}
}

func TestHandleFile_SortTestdataScores(t *testing.T) {
	t.Parallel()

	fd, err := os.Open(filepath.Join("..", "..", "testdata", "scores.csv"))
	require.NoError(t, err)
	t.Cleanup(func() { fd.Close() })

	got, err := HandleFile(FileOpts{
		Delimiter: ',',
		KeepCols:  []string{"score", "name"},
		SortKeys:  []SortKey{{Name: "score", Desc: true}},
	}, fd)
	require.NoError(t, err)
	require.Equal(t, []Line{
		{"45", "b"},
		{"32", "c"},
		{"12", "a"},
	}, got.Lines)
}

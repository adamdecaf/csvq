package main

import (
	"testing"

	"github.com/adamdecaf/csvq/internal/cli"

	"github.com/stretchr/testify/require"
)

func TestSplitStringList(t *testing.T) {
	got := splitStringList("")
	require.Empty(t, got)
}

func TestParseSortKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want []cli.SortKey
	}{
		{
			name: "empty",
			args: []string{"file.csv"},
			want: nil,
		},
		{
			name: "desc then asc keeps invocation order",
			args: []string{"-sort.dsc", "score", "-sort.asc", "name", "file.csv"},
			want: []cli.SortKey{
				{Name: "score", Desc: true},
				{Name: "name"},
			},
		},
		{
			name: "equals form and comma list",
			args: []string{"-sort.asc=name,date", "-sort.dsc=score"},
			want: []cli.SortKey{
				{Name: "name"},
				{Name: "date"},
				{Name: "score", Desc: true},
			},
		},
		{
			name: "double dash",
			args: []string{"--sort.dsc", "score", "--sort.asc=name"},
			want: []cli.SortKey{
				{Name: "score", Desc: true},
				{Name: "name"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, parseSortKeys(tt.args))
		})
	}
}

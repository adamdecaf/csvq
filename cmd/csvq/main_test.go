package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitStringList(t *testing.T) {
	got := splitStringList("")
	require.Empty(t, got)
}

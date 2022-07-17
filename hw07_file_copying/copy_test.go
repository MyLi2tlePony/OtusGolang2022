package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("offset 0, limit 0", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 0))
	})

	t.Run("offset 0, limit 10", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 10))
	})

	t.Run("offset 0, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 1000))
	})
}

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")

		require.NotNil(t, env)
		require.Nil(t, err)

		require.Equal(t, EnvValue{
			"bar",
			false,
		}, env["BAR"])

		require.Equal(t, EnvValue{
			"",
			false,
		}, env["EMPTY"])

		require.Equal(t, EnvValue{
			"   foo\nwith new line",
			false,
		}, env["FOO"])

		require.Equal(t, EnvValue{
			"\"hello\"",
			false,
		}, env["HELLO"])

		require.Equal(t, EnvValue{
			"",
			true,
		}, env["UNSET"])

		require.Equal(t, EnvValue{
			"TAB",
			false,
		}, env["TAB"])

		require.Equal(t, EnvValue{
			"WHITESPACE",
			false,
		}, env["WHITESPACE"])

		require.Equal(t, EnvValue{
			"WHITESPACEANDTAB",
			false,
		}, env["WHITESPACEANDTAB"])
	})
}

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	input := "testdata/input.txt"
	output := t.TempDir() + "/out.txt"

	t.Run("identical in and output files", func(t *testing.T) {
		res := Copy(input, "testdata/../"+input, 100000, 0)
		require.Error(t, res)
	})

	t.Run("offset less zero", func(t *testing.T) {
		res := Copy(input, output, -1, 0)
		require.Error(t, res)
	})

	t.Run("limit less zero", func(t *testing.T) {
		res := Copy(input, output, 0, -1)
		require.Error(t, res)
	})

	t.Run("offset over file size", func(t *testing.T) {
		res := Copy(input, output, 100000, 0)
		require.Error(t, res)
	})

	t.Run("limit over file size", func(t *testing.T) {
		res := Copy(input, output, 0, 1000000000)
		require.Nil(t, res)
		require.FileExists(t, input)
		require.FileExists(t, output)

		inputStat, inputErr := os.Stat(input)
		outputStat, outputErr := os.Stat(output)
		require.Nil(t, inputErr)
		require.Nil(t, outputErr)
		require.Equal(t, inputStat.Size(), outputStat.Size())
	})
}

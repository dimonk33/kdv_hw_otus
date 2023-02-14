package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func removeFile(path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		fmt.Printf("error delete file %s\r\n", path)
	}
}

func TestCopy(t *testing.T) {
	t.Run("offset over file size", func(t *testing.T) {
		input := "testdata/input.txt"
		output := "testdata/out.txt"
		res := Copy(input, output, 100000, 0)
		defer removeFile(output)
		require.Error(t, res)
	})

	t.Run("limit over file size", func(t *testing.T) {
		input := "testdata/input.txt"
		output := "testdata/out.txt"

		res := Copy(input, output, 0, 1000000000)
		defer removeFile(output)
		require.Nil(t, res)

		inputStat, inputErr := os.Stat(input)
		outputStat, outputErr := os.Stat(output)
		require.Nil(t, inputErr)
		require.Nil(t, outputErr)
		require.Equal(t, inputStat.Size(), outputStat.Size())
	})
}

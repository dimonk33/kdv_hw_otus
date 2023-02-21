package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("set environment", func(t *testing.T) {
		err := os.Setenv("TEST2", "Test2")
		require.Nil(t, err)

		testEnvironment := Environment{
			"TEST1": EnvValue{Value: "Test1", NeedRemove: false},
			"TEST2": EnvValue{Value: "", NeedRemove: true},
		}
		setEnvironment(testEnvironment)

		test1, err1 := os.LookupEnv("TEST1")
		require.True(t, err1)
		require.Equal(t, test1, "Test1")

		_, err2 := os.LookupEnv("TEST2")
		require.False(t, err2)
	})

	t.Run("exec cmd", func(t *testing.T) {
		res := execCmd("bash", []string{"-c", "echo", "1"})
		require.Equal(t, 0, res)
	})
}

package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("filter env key", func(t *testing.T) {
		res := filterEnvKey("test= ; \t")
		require.Equal(t, res, "test")
	})

	t.Run("filter env value", func(t *testing.T) {
		res := filterEnvValue("test=\x00\x00    \n test2\n")
		require.Equal(t, "test=\n\n", res)
	})

	t.Run("get env from dir", func(t *testing.T) {
		res, err := ReadDir("./testdata/env")
		require.Nil(t, err)
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		require.Equal(t, expected, res)
	})

	t.Run("bad dir path", func(t *testing.T) {
		_, err := ReadDir("../testdata/env")
		require.NotNil(t, err)
	})

	t.Run("skip env with =", func(t *testing.T) {
		testDir := "./testdata/env"
		testFile := testDir + "/" + "test="
		file, err := os.Create(testFile)
		if err != nil {
			log.Fatal(err)
		}
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := os.Remove(testFile)
			if err != nil {
				log.Fatal(err)
			}
		}()
		res, errEnv := ReadDir(testDir)
		require.Nil(t, errEnv)
		_, found := res["test="]
		require.False(t, found)
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		require.Equal(t, expected, res)
	})
}

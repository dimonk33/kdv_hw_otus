package main

import (
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesEnv := make(Environment)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, errInfo := file.Info()
		if errInfo != nil {
			continue
		}

		key := filterEnvKey(file.Name())
		if info.Size() == 0 {
			filesEnv[key] = EnvValue{Value: "", NeedRemove: true}
			continue
		}
		content, errRead := os.ReadFile(path.Join(dir, file.Name()))
		if errRead != nil {
			continue
		}
		filesEnv[key] = EnvValue{Value: filterEnvValue(string(content)), NeedRemove: false}
	}

	return filesEnv, nil
}

func filterEnvKey(key string) string {
	return strings.TrimRight(key, "=; \t")
}

func filterEnvValue(value string) string {
	lines := strings.Split(value, "\n")
	if len(lines) == 0 {
		return ""
	}
	return strings.TrimRight(strings.ReplaceAll(lines[0], "\x00", "\n"), " ")
}

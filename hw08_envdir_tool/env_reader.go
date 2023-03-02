package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	dirAbsPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dirFiles, err := os.ReadDir(dirAbsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range dirFiles {
		var envValue EnvValue

		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		data, err := os.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			envValue.NeedRemove = true
		} else {
			envValue.Value = clean(data)
		}

		env[fileName] = envValue
	}

	return env, err
}

func clean(v []byte) string {
	v = bytes.Split(v, []byte("\n"))[0]

	return strings.ReplaceAll(strings.TrimRight(string(v), " \t\n"), "\x00", "\n")
}

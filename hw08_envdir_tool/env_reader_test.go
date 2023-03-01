package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirNotFound(t *testing.T) {
	dirPath := filepath.Join(os.TempDir(), "/testdata")
	_, err := ReadDir(dirPath)

	var errPath *os.PathError

	require.True(t, errors.As(err, &errPath))
}

func TestReadDir(t *testing.T) {
	testCases := map[string]string{
		"CASE_DELETE":     "",
		"CASE_END":        "test \t",
		"CASE_REPLACE":    string([]byte{'t', 'e', 's', 't', '\x00', 't', 'e', 's', 't'}),
		"CASE_FIRST_LINE": "first\ntwo\nthree\nfour",
		"CASE_POSITIVE_1": "value1",
		"CASE_POSITIVE_2": "VALUE2",
	}

	expected := Environment{
		"CASE_DELETE": EnvValue{
			NeedRemove: true,
		},
		"CASE_END": EnvValue{
			Value: "test",
		},
		"CASE_REPLACE": EnvValue{
			Value: "test\ntest",
		},
		"CASE_FIRST_LINE": EnvValue{
			Value: "first",
		},
		"CASE_POSITIVE_1": EnvValue{
			Value: "value1",
		},
		"CASE_POSITIVE_2": EnvValue{
			Value: "VALUE2",
		},
	}

	tmpDirPath := tmpTestDirPath(t, testCases)
	result, err := ReadDir(tmpDirPath)

	require.NoError(t, err)
	require.Equal(t, expected, result)

	err = os.RemoveAll(tmpDirPath)

	require.NoError(t, err)
}

func tmpTestDirPath(t *testing.T, testCases map[string]string) string {
	t.Helper()

	var err error

	testDirPath := filepath.Join(os.TempDir(), "/hw08_envdir_tool")

	_ = os.RemoveAll(testDirPath)

	err = os.Mkdir(testDirPath, os.ModePerm)
	require.NoError(t, err)

	for fileName, value := range testCases {
		fileName = filepath.Join(testDirPath, "/"+fileName)
		err := os.WriteFile(fileName, []byte(value), os.ModePerm)

		require.NoError(t, err)
	}

	return testDirPath
}

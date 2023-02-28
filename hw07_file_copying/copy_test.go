package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var (
		fromPath, toPath string
		limit, offset    int64
	)

	t.Run("File not found", func(t *testing.T) {
		fromPath = "./testdata/none.txt"
		toPath = "./test1.txt"

		result := Copy(fromPath, toPath, offset, limit)

		require.Truef(t, errors.Is(result, ErrUnsupportedFile), "%v", result)

		_ = os.Remove(toPath)
	})

	t.Run("Unsupported file", func(t *testing.T) {
		fromPath = "/dev/urandom"
		toPath = "./test2.txt"

		result := Copy(fromPath, toPath, offset, limit)

		require.Truef(t, errors.Is(result, ErrUnsupportedFile), "%v", result)

		_ = os.Remove(toPath)
	})

	t.Run("Offset exceeds file size", func(t *testing.T) {
		fromPath = "./testdata/offset_exceed.txt"
		toPath = "./test3.txt"
		offset = 5000

		result := Copy(fromPath, toPath, offset, limit)

		require.Truef(t, errors.Is(result, ErrOffsetExceedsFileSize), "%v", result)

		_ = os.Remove(toPath)
	})

	t.Run("Success copy file", func(t *testing.T) {
		fromPath = "./testdata/input.txt"
		toPath = "./test4.txt"

		result := Copy(fromPath, toPath, offset, limit)

		require.NoError(t, result)
		require.Nil(t, result)

		_ = os.Remove(toPath)
	})
}

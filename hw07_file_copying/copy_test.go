package main

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("offset 0, limit 0", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 0))
		require.Equal(t, true, FilesIsEqual("out.txt", "testdata/out_offset0_limit0.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 0, limit 10", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 10))
		require.Equal(t, true, FilesIsEqual("out.txt", "testdata/out_offset0_limit10.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 0, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 1000))
		require.Equal(t, true, FilesIsEqual("out.txt", "testdata/out_offset0_limit1000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 0, limit 10000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 10000))
		require.Equal(t, true, FilesIsEqual("out.txt", "testdata/out_offset0_limit10000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 100, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 100, 1000))
		require.Equal(t, true, FilesIsEqual("out.txt", "testdata/out_offset100_limit1000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 6000, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 6000, 1000))
		require.Equal(t, true, FilesIsEqual("out.txt", "testdata/out_offset6000_limit1000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})
}

func FilesIsEqual(filePath1, filePath2 string) (result bool) {
	file1, err := os.Open(filePath1)
	if err != nil {
		return false
	}

	defer func() {
		closeErr := file1.Close()
		if closeErr != nil {
			result = false
		}
	}()

	file2, err := os.Open(filePath2)
	if err != nil {
		return false
	}

	defer func() {
		closeErr := file2.Close()
		if closeErr != nil {
			result = false
		}
	}()

	bytesNumber := int64(100)
	for i := int64(0); ; i++ {
		bytes1 := make([]byte, bytesNumber)
		n1, err1 := file1.ReadAt(bytes1, bytesNumber*i)

		bytes2 := make([]byte, bytesNumber)
		n2, err2 := file2.ReadAt(bytes2, bytesNumber*i)

		if !bytes.Equal(bytes1, bytes2) {
			return false
		}

		if !errors.Is(err1, err2) {
			return false
		}

		if n2 == 0 && n1 == 0 {
			break
		}
	}

	return true
}

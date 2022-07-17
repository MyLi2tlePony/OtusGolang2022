package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile          = errors.New("unsupported file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrOffsetIsGreaterThanLimit = errors.New("offset is greater than limit")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fromFileSize := fromFileInfo.Size()

	if !fromFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if limit == 0 {
		limit = fromFileSize
	}
	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}
	if offset > limit {
		return ErrOffsetIsGreaterThanLimit
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := fromFile.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	if _, err = fromFile.Seek(offset, 0); err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := toFile.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	bar := StartNewProgressBar(limit)

	for offset < limit {
		written, err := io.CopyN(toFile, fromFile, 1)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		offset += written
		bar.Add(written)
	}

	return nil
}

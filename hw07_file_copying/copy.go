package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)

	switch {
	case err != nil || fileInfo.IsDir() || fileInfo.Size() == 0:
		return ErrUnsupportedFile

	case offset > fileInfo.Size():
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)

	if err != nil {
		return err
	}

	defer func() {
		if err = fromFile.Close(); err != nil {
			log.Panicf("Failed to close file: %v", err)
		}
	}()

	seekPosition, err := fromFile.Seek(offset, io.SeekStart)

	if err != nil {
		return err
	}

	limitBar := fileInfo.Size() - seekPosition
	limitBuffer := limit

	if limit < limitBar {
		limitBar = limit
	}

	if limit == 0 {
		limitBuffer = fileInfo.Size()
	}

	bar := pb.Simple.Start64(limitBar)

	defer bar.Finish()

	reader := io.LimitReader(fromFile, limitBuffer)
	readerBar := bar.NewProxyReader(reader)

	toFile, err := os.Create(toPath)

	if err != nil {
		return err
	}

	if _, err = io.Copy(toFile, readerBar); err != nil {
		return err
	}

	return nil
}

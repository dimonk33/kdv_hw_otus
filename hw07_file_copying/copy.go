package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFileNotExist          = errors.New("file not exist")
	ErrCreateFile            = errors.New("file create error")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSeekFile              = errors.New("file seek error")
	ErrCopyFile              = errors.New("file copy error")
	ErrPathFiles             = errors.New("path for files error")
	ErrParams                = errors.New("limit or offset error")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromPath = path.Clean(fromPath)
	toPath = path.Clean(toPath)

	if fromPath == toPath {
		return ErrPathFiles
	}

	if offset < 0 || limit < 0 {
		return ErrParams
	}

	fileFrom, errFrom := os.Open(fromPath)
	if errFrom != nil {
		if os.IsNotExist(errFrom) {
			return ErrFileNotExist
		}
		return errFrom
	}
	defer closeFile(fileFrom)

	info, errSize := fileFrom.Stat()
	fileSize := info.Size()
	if errSize != nil {
		return ErrUnsupportedFile
	}

	if offset >= fileSize {
		return ErrOffsetExceedsFileSize
	}

	_ = os.Remove(toPath)
	fileTo, errTo := os.Create(toPath)
	if errTo != nil {
		return ErrCreateFile
	}
	defer closeFile(fileTo)

	copySize := fileSize - offset
	if limit > 0 && limit < copySize {
		copySize = limit
	}
	_, errSeek := fileFrom.Seek(offset, io.SeekStart)
	if errSeek != nil {
		return ErrSeekFile
	}

	bar := pb.Full.Start64(copySize)
	barReader := bar.NewProxyReader(fileFrom)

	_, errCopy := io.CopyN(fileTo, barReader, copySize)

	bar.Finish()
	fmt.Println("")

	if errCopy != nil {
		return ErrCopyFile
	}

	return nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Printf("%v", err)
	}
}

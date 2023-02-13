package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrFileNotExist = errors.New("file not exist")
	ErrCreateFile   = errors.New("file create error")
	//ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSeekFile              = errors.New("file seek error")
	ErrCopyFile              = errors.New("file copy error")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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
		return errSize
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
	_, errSeek := fileFrom.Seek(offset, 0)
	if errSeek != nil {
		return ErrSeekFile
	}
	_, errCopy := io.CopyN(fileTo, fileFrom, copySize)
	if errCopy != nil {
		return ErrCopyFile
	}

	return nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Printf("%v", err)
	}
}

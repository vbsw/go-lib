/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package fs provides various file function.
package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrBufferOutOfMemory = errors.New("buffer out of memory")
	errDummy             = errors.New("dummy")
)

type File struct {
	Err  error
	Info os.FileInfo
}

// FileReader reads files into buffer.
type FileReader struct {
	File
	Offset int64
	file   *os.File
	Buffer []byte
	NRead  int
}

// Open opens file for reading and initializes Info.
// It returns true if file was successfully opened.
// Error is stored in Err.
func (reader *FileReader) Open(path string) bool {
	reader.Offset, reader.NRead = 0, 0
	reader.file, reader.Err = os.Open(path)
	if reader.Err == nil {
		reader.Info, reader.Err = reader.file.Stat()
		if reader.Err != nil {
			reader.file.Close()
		}
	}
	if reader.Err != nil {
		reader.file = nil
	}
	return reader.file != nil
}

// Read copies the last keepN bytes to the beginning
// of the buffer and then reads from the file into the buffer.
// Buffer is starting at offset keepN. It returns true if any
// bytes have been read and no error encountered.
// The error is stored in Err unless it is io.EOF.
func (reader *FileReader) Read(keepN int) bool {
	var err error
	if keepN <= 0 {
		reader.NRead, err = reader.file.Read(reader.Buffer)
	} else if keepN < len(reader.Buffer) {
		copy(reader.Buffer, reader.Buffer[len(reader.Buffer)-keepN:])
		reader.NRead, err = reader.file.Read(reader.Buffer[keepN:])
		reader.NRead += keepN
	} else {
		err = ErrBufferOutOfMemory
	}
	if err == io.EOF {
		reader.Err = nil
	} else {
		reader.Err = err
	}
	reader.Offset += int64(reader.NRead)
	return reader.NRead > 0 && err == nil
}

// IsOpen returns true if file is open.
func (reader *FileReader) IsOpen() bool {
	return reader.file != nil
}

// Seek sets the offset for the next Read on file.
// It returns true if seek was successful.
// Error is stored in Err.
func (reader *FileReader) Seek(offset int64) bool {
	reader.Offset, reader.Err = reader.file.Seek(offset, io.SeekStart)
	return reader.Err == nil
}

// Close closes the file.
// Err is set only if it was previously nil.
func (reader *FileReader) Close() {
	if reader.file != nil {
		if reader.Err == nil {
			reader.Err = reader.file.Close()
		} else {
			reader.file.Close()
		}
		reader.file = nil
	}
}

// Stat calls os.Stat(path), stores result in Info.
// It returns true if file exists.
// Error is stored in Err.
func (file *File) Stat(path string) bool {
	file.Info, file.Err = os.Stat(path)
	return file.Info != nil && (file.Err == nil || !os.IsNotExist(file.Err))
}

// IsDir returns true if path is a directory.
// Error is stored in Err.
func (file *File) IsDir(path string) bool {
	file.Info, file.Err = os.Stat(path)
	if file.Err == nil && file.Info != nil {
		return file.Info.IsDir()
	}
	return false
}

// IsRegular returns true if path is a regular file.
// Error is stored in Err.
func (file *File) IsRegular(path string) bool {
	file.Info, file.Err = os.Stat(path)
	if file.Err == nil && file.Info != nil {
		return file.Info.Mode().IsRegular()
	}
	return false
}

// IsEmpty returns whether directory or file is empty.
// A directory is empty if it has only empty directories and empty files.
// A file is empty if it is a regular file with size 0.
func (file *File) IsEmpty(path string) bool {
	file.Info, file.Err = os.Stat(path)
	if file.Info != nil {
		if file.Info.Mode().IsRegular() {
			if file.Err == nil {
				return file.Info.Size() == 0
			}
			return os.IsNotExist(file.Err)
		} else if file.Info.IsDir() {
			return isDirEmpty(file, path)
		}
	}
	return false
}

func isDirEmpty(file *File, path string) bool {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info != nil {
				if info.IsDir() || info.Mode().IsRegular() && info.Size() == 0 {
					return nil
				}
				return errDummy
			}
			return errDummy
		} else if os.IsNotExist(file.Err) {
			return nil
		}
		return errDummy
	})
	return err == nil
}

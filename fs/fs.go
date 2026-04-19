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
		err = errors.New("buffer out of memory")
	}
	if err == io.EOF {
		reader.Err = nil
	} else {
		reader.Err = err
	}
	reader.Offset += int64(reader.NRead)
	return reader.NRead > 0 && reader.Err == nil
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

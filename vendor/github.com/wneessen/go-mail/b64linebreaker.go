// SPDX-FileCopyrightText: The go-mail Authors
//
// SPDX-License-Identifier: MIT

package mail

import (
	"errors"
	"io"
)

// newlineBytes is a byte slice representation of the SingleNewLine constant used for line breaking
// in encoding processes.
var newlineBytes = []byte(SingleNewLine)

// ErrNoOutWriter is the error returned when no io.Writer is set for Base64LineBreaker.
var ErrNoOutWriter = errors.New("no io.Writer set for Base64LineBreaker")

// Base64LineBreaker handles base64 encoding with the insertion of new lines after a certain number
// of characters.
//
// This struct is used to manage base64 encoding while ensuring that new lines are inserted after
// reaching a specific line length. It satisfies the io.WriteCloser interface.
//
// References:
//   - https://datatracker.ietf.org/doc/html/rfc2045 (Base64 and line length limitations)
type Base64LineBreaker struct {
	line [MaxBodyLength]byte
	used int
	out  io.Writer
}

// Write writes data to the Base64LineBreaker, ensuring lines do not exceed MaxBodyLength.
//
// This method writes the provided data to the Base64LineBreaker. It ensures that the written
// lines do not exceed the MaxBodyLength. If the data exceeds the limit, it handles the
// continuation by splitting the data and writing new lines as necessary.
//
// Parameters:
//   - data: A byte slice containing the data to be written.
//
// Returns:
//   - numBytes: The number of bytes written.
//   - err: An error if one occurred during the write operation.
func (l *Base64LineBreaker) Write(data []byte) (numBytes int, err error) {
	if l.out == nil {
		err = ErrNoOutWriter
		return
	}
	if l.used+len(data) < MaxBodyLength {
		copy(l.line[l.used:], data)
		l.used += len(data)
		return len(data), nil
	}

	numBytes, err = l.out.Write(l.line[0:l.used])
	if err != nil {
		return
	}
	excess := MaxBodyLength - l.used
	l.used = 0

	numBytes, err = l.out.Write(data[0:excess])
	if err != nil {
		return
	}

	numBytes, err = l.out.Write(newlineBytes)
	if err != nil {
		return
	}

	return l.Write(data[excess:])
}

// Close finalizes the Base64LineBreaker, writing any remaining buffered data and appending a newline.
//
// This method ensures that any remaining data in the buffer is written to the output and appends
// a newline. It is used to finalize the Base64LineBreaker and should be called when no more data
// is expected to be written.
//
// Returns:
//   - err: An error if one occurred during the final write operation.
func (l *Base64LineBreaker) Close() (err error) {
	if l.used > 0 {
		_, err = l.out.Write(l.line[0:l.used])
		if err != nil {
			return
		}
		_, err = l.out.Write(newlineBytes)
	}

	return
}

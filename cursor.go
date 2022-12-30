// Copyright (c) 2023 Steven Stallion <sstallion@gmail.com>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package smclcd

import (
	"errors"
	"io"
)

const (
	Lines   = 2
	Columns = 16
)

type cursor int

func (pos *cursor) Advance(n int) error {
	*pos = cursor(int(*pos) + n)
	return pos.Error()
}

func (pos *cursor) Move(y, x int) error {
	*pos = cursor(y*Columns + x)
	return pos.Error()
}

func (pos *cursor) Error() error {
	switch {
	case *pos < 0:
		return errors.New("cursor: negative position")
	case *pos >= Lines*Columns:
		return io.EOF
	default:
		return nil
	}
}

func (pos *cursor) Remaining() int {
	return Columns - int(*pos%Columns)
}

func (pos *cursor) Byte() byte {
	return pCursorPos + pCursorLn*byte(*pos/Columns) + byte(*pos%Columns)
}

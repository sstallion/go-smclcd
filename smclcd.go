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
	"bytes"
	"fmt"
	"io"
	"unicode"

	"github.com/sstallion/go-hid"
	"github.com/sstallion/go-tools/util"
)

const (
	VendorID  = 0x15d9 // SMC
	ProductID = 0x1133 // SuperMicro LCD Display
)

type Backlight byte

//go:generate stringer -type Backlight -trimprefix=Backlight

const (
	BacklightOff Backlight = iota
	BacklightOn
)

type Cursor byte

//go:generate stringer -type Cursor -trimprefix=Cursor

const (
	CursorOff Cursor = iota
	CursorBlock
	CursorUnderline
	CursorBoth
)

type Key struct {
	Code  KeyCode
	Event KeyEvent
}

type KeyCode byte

//go:generate stringer -type KeyCode -trimprefix=Key

const (
	KeyUp KeyCode = iota
	KeyRight
	KeyLeft
	KeyDown
	KeyEnter
	KeyCancel
)

type KeyEvent byte

//go:generate stringer -type KeyEvent -trimprefix=Key

const (
	KeyRelease KeyEvent = iota
	KeyPress
)

type LCD struct {
	device *hid.Device
	pos    cursor
}

func Open(serial string) (l *LCD, err error) {
	var device *hid.Device
	if device, err = hid.Open(VendorID, ProductID, serial); err != nil {
		return
	}
	l = &LCD{device: device}
	return
}

func OpenFirst() (l *LCD, err error) {
	var device *hid.Device
	if device, err = hid.OpenFirst(VendorID, ProductID); err != nil {
		return
	}
	l = &LCD{device: device}
	return
}

func OpenPath(path string) (l *LCD, err error) {
	var device *hid.Device
	if device, err = hid.OpenPath(path); err != nil {
		return
	}
	l = &LCD{device: device}
	return
}

func (l *LCD) Close() error {
	return l.device.Close()
}

func (l *LCD) recvInputReport(p, prefix []byte) (err error) {
	prefix = append([]byte{inputReportID}, prefix...)
	b := make([]byte, inputReportLen)
	for {
		if _, err = l.device.Read(b); err != nil {
			return
		}
		logReport(b)
		if bytes.HasPrefix(b, prefix) {
			break
		}
	}
	copy(p, bytes.TrimPrefix(b[:len(b)-1], prefix))
	return
}

func (l *LCD) sendOutputReport(p []byte) (err error) {
	b := make([]byte, outputReportLen)
	b[0] = outputReportID
	copy(b[1:len(b)-1], p)
	b[len(b)-1] = checksum(b)
	if _, err = l.device.Write(b); err != nil {
		return
	}
	logReport(b)
	return
}

func (l *LCD) Version() (s string, err error) {
	prefix := []byte{pVersion}
	if err = l.sendOutputReport(prefix); err != nil {
		return
	}

	b := make([]byte, inputReportCmdLen)
	if err = l.recvInputReport(b, prefix); err != nil {
		return
	}
	s = fmt.Sprintf("%x.%x", b[0], b[1])
	return
}

func (l *LCD) Clear() error {
	l.pos.Move(0, 0)
	b := []byte{pLCD, pControl, pClear}
	return l.sendOutputReport(b)
}

func (l *LCD) Home() error {
	l.pos.Move(0, 0)
	b := []byte{pLCD, pControl, pHome}
	return l.sendOutputReport(b)
}

func (l *LCD) SetCursor(state Cursor) error {
	b := []byte{pLCD, pControl, pCursor + byte(state)}
	return l.sendOutputReport(b)
}

func (l *LCD) AdvanceCursor(n int) (err error) {
	if err = l.pos.Advance(n); err != nil {
		return
	}
	b := []byte{pLCD, pControl, l.pos.Byte()}
	return l.sendOutputReport(b)
}

func (l *LCD) MoveCursor(y, x int) (err error) {
	if err = l.pos.Move(y, x); err != nil {
		return
	}
	b := []byte{pLCD, pControl, l.pos.Byte()}
	return l.sendOutputReport(b)
}

func (l *LCD) Write(p []byte) (n int, err error) {
	lines := bytes.Split(p, []byte("\n"))
	for i, line := range lines {
		var m int
		var remaining = l.pos.Remaining()
		m, err = l.writeRaw(line)
		if n += m; err != nil {
			return
		}

		if i+1 < len(lines) && len(line) < remaining {
			b := bytes.Repeat([]byte(" "), l.pos.Remaining())
			if _, err = l.writeRaw(b); err != nil {
				return
			}
			n++
		}
	}
	return
}

func (l *LCD) writeRaw(p []byte) (n int, err error) {
	prefix := []byte{pLCD, pWrite}
	p = bytes.Map(func(r rune) rune {
		if !unicode.IsPrint(r) {
			return '?'
		}
		return r
	}, p)
	for n < len(p) {
		if err = l.pos.Error(); err != nil {
			return
		}

		m := util.Min(len(p[n:]), util.Min(l.pos.Remaining(), outputReportDataLen))
		b := append(prefix, p[n:n+m]...)
		if err = l.sendOutputReport(b); err != nil {
			return
		}
		n += m

		if err = l.AdvanceCursor(m); err != nil {
			return
		}
	}
	return
}

func (l *LCD) Read(p []byte) (n int, err error) {
	for n < len(p) {
		var b = make([]byte, l.pos.Remaining())
		var m int
		if m, err = l.readRaw(b); err != nil {
			if err != io.EOF {
				return
			}
		}
		n += copy(p[n:], b)

		if n < len(p) && m == len(b) {
			p[n] = '\n'
			n++
		}
	}
	return
}

func (l *LCD) readRaw(p []byte) (n int, err error) {
	prefix := []byte{pLCD, pRead}
	for n < len(p) {
		if err = l.pos.Error(); err != nil {
			return
		}

		if err = l.sendOutputReport(prefix); err != nil {
			return
		}

		m := util.Min(len(p[n:]), util.Min(l.pos.Remaining(), inputReportDataLen))
		if err = l.recvInputReport(p[n:n+m], prefix); err != nil {
			return
		}
		n += m

		if err = l.AdvanceCursor(m); err != nil {
			return
		}
	}
	return
}

func (l *LCD) GetInput() (key Key, err error) {
	prefix := []byte{pKeyInput}
	b := make([]byte, inputReportDataLen)
	if err = l.recvInputReport(b, prefix); err != nil {
		return
	}
	key = Key{
		Code:  KeyCode(b[0]),
		Event: KeyEvent(b[1]),
	}
	return
}

func (l *LCD) SetBacklight(state Backlight) error {
	b := []byte{pBacklight, byte(state)}
	return l.sendOutputReport(b)
}

func (l *LCD) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(l, a...)
}

func (l *LCD) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(l, format, a...)
}

func (l *LCD) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(l, a...)
}

func (l *LCD) Scan(a ...interface{}) (int, error) {
	return fmt.Fscan(l, a...)
}

func (l *LCD) Scanf(format string, a ...interface{}) (int, error) {
	return fmt.Fscanf(l, format, a...)
}

func (l *LCD) Scanln(a ...interface{}) (int, error) {
	return fmt.Fscanln(l, a...)
}

var _ io.ReadWriteCloser = (*LCD)(nil)

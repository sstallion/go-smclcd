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

// ## HID Reports
//
// First byte is report ID; AA for input BB for output
// Last byte appears to be checksum; 2s complement
// Both reports are 16 bytes long
// chksum choice seems odd; interrupt endpoints are guaranteed delivery unsure of
// why this was added.
// report followed by one command byte, and an optional subcommand
//
// 1.  LCD (`02`)
//
//  1. Control (`00`)
//
//     00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |BB|02|00|01|00|00|00|00|00|00|00|00|00|00|00|42|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |ID|     |XX|            RESERVED            |CK|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
//     ID HID Report ID (BB)
//     XX Command:
//     01 Clear
//     02 Home
//     0C Cursor Off
//     0D Block Cursor*
//     0E Underline Cursor
//     0F Block & Underline Cursor
//     8X Move Cursor to Line 0 (+ Column)
//     CX Move Cursor to Line 1 (+ Column)
//     CK Checksum
//
//     * Might be unintentional; decompiled java seems to indicate this
//     is not a valid option.
//
//  2. Version (01)
//
//     TODO add version section and remove markdown from docs
//
//  2. Write (`02`)
//
//     00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |BB|02|02|00|00|00|00|00|00|00|00|00|00|00|00|33|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |ID|     |               DATA                |CK|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
//     ID HID Report ID (BB)
//     CK Checksum
//
//  3. Read (`03`)
//
//     00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |BB|02|03|00|00|00|00|00|00|00|00|00|00|00|00|40|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |ID|     |             RESERVED              |CK|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
//     ID HID Report ID (BB)
//     CK Checksum
//
//     00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |AA|02|03|20|20|20|20|20|20|20|20|20|20|20|20|D1|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |ID|     |               DATA                |CK|
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
//     ID HID Report ID (AA)
//     CK Checksum
//
// 2.  Key Input (`03`)
//
//	 00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|AA|03|00|00|00|00|00|00|00|00|00|00|00|00|00|53|
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|ID|  |XX|YY|            RESERVED            |CK|
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
//	 ID HID Report ID (AA)
//	 XX Key Code:
//	    00 Up
//	    01 Right
//	    02 Left
//	    03 Down
//	    04 Enter
//	    05 Cancel
//	 YY Key Event:
//	    00 Key Pressed
//	    01 Key Released
//	 CK Checksum
//
// 3.  Set Backlight (`07`)
//
//	 00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|BB|07|00|00|00|00|00|00|00|00|00|00|00|00|00|3E|
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//	|ID|  |XX|             RESERVED              |CK|
//	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
//	 ID HID Report ID (BB)
//	 XX State:
//	    00 Off
//	    01 On
//	 CK Checksum
const (
	inputReportID      = 0xaa
	inputReportLen     = 16
	inputReportCmdLen  = inputReportLen - 2
	inputReportDataLen = inputReportCmdLen - 2

	outputReportID      = 0xbb
	outputReportLen     = 16
	outputReportCmdLen  = outputReportLen - 2
	outputReportDataLen = outputReportCmdLen - 2

	// Command bytes
	pVersion   = 0x01
	pLCD       = 0x02
	pKeyInput  = 0x03
	pBacklight = 0x07

	// LCD bytes
	pControl = 0x00
	pWrite   = 0x02
	pRead    = 0x03

	// Control bytes
	pClear     = 0x01
	pHome      = 0x02
	pCursor    = 0x0c
	pCursorPos = 0x80
	pCursorLn  = 0x40
)

func checksum(b []byte) byte {
	var val byte
	for _, v := range b {
		val += v
	}
	return ^val + 1
}

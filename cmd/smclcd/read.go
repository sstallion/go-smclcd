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

package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/sstallion/go-smclcd"
	"github.com/sstallion/go-tools/command"
)

type readCmd struct {
	flags   *flag.FlagSet
	n, x, y uint
}

func init() {
	cmd := &readCmd{flags: flag.NewFlagSet("read", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	cmd.flags.UintVar(&cmd.n, "n", 0, "`count`")
	cmd.flags.UintVar(&cmd.x, "x", 0, "`col`")
	cmd.flags.UintVar(&cmd.y, "y", 0, "`line`")
	command.Add(cmd)
}

func (cmd *readCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *readCmd) Description() string {
	return "Read from display"
}

func (cmd *readCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }} [-n count] [-x col] [-y line]

Flags:

  {{ call .PrintDefaults }}

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *readCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	if cmd.n == 0 { // default to entire display
		cmd.n = smclcd.Lines * (smclcd.Columns + 1)
	}
	args := cmd.flags.Args()
	if len(args) != 0 {
		return command.ErrNArg
	}
	return nil
}

func (cmd *readCmd) Run() error {
	l, err := openLCD()
	if err != nil {
		return err
	}
	defer l.Close()

	y, x := int(cmd.y), int(cmd.x)
	if err = l.MoveCursor(y, x); err != nil {
		return err
	}

	b := make([]byte, cmd.n)
	if _, err = l.Read(b); err != nil {
		if err != io.EOF {
			return err
		}
	}
	fmt.Println(strings.TrimRight(string(b), "\n"))
	return nil
}

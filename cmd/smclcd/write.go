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
	"io"
	"strings"

	"github.com/sstallion/go-tools/command"
)

type writeCmd struct {
	flags *flag.FlagSet
	clear bool
	x, y  uint
	args  []string
}

func init() {
	cmd := &writeCmd{flags: flag.NewFlagSet("write", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	cmd.flags.BoolVar(&cmd.clear, "clear", false, "TODO")
	cmd.flags.UintVar(&cmd.x, "x", 0, "`col`")
	cmd.flags.UintVar(&cmd.y, "y", 0, "`line`")
	command.Add(cmd)
}

func (cmd *writeCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *writeCmd) Description() string {
	return "Write to display"
}

func (cmd *writeCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }} [-x col] [-y line] arguments...

Flags:

  {{ call .PrintDefaults }}

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *writeCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	args := cmd.flags.Args()
	if len(args) < 1 {
		return command.ErrNArg
	}
	cmd.args = args
	return nil
}

func (cmd *writeCmd) Run() error {
	l, err := openLCD()
	if err != nil {
		return err
	}
	defer l.Close()

	if cmd.clear {
		if err = l.Clear(); err != nil {
			return err
		}
	}

	y, x := int(cmd.y), int(cmd.x)
	if err = l.MoveCursor(y, x); err != nil {
		return err
	}

	b := []byte(strings.Join(cmd.args, " "))
	if _, err = l.Write(b); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}
